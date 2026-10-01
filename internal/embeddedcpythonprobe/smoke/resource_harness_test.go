package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	python "github.com/goccy/go-python"
)

const childDeadline = 30 * time.Second
const cancelDelay = 10 * time.Millisecond

type resourceCase struct {
	Kind string
	N    int
}

func resourceCases() []resourceCase {
	return []resourceCase{{"depth", 100}, {"depth", 500}, {"broad", 10000}, {"broad", 100000}, {"token", 1 << 20}, {"token", 8 << 20}, {"matches", 10000}, {"matches", 100000}, {"files", 100}, {"files", 1000}, {"invalid-last", 1}, {"cancel-native", 100000}, {"cancel-walk", 100000}}
}

// Reject parameters BEFORE allocation. This is an experiment input ceiling,
// not a production parser/memory limit.
func resourceInputs(c resourceCase) ([]probeInput, error) {
	limits := map[string]int{"depth": 1000, "broad": 100000, "token": 8 << 20, "matches": 100000, "files": 1000, "invalid-last": 1, "cancel-native": 100000, "cancel-walk": 100000}
	limit, ok := limits[c.Kind]
	if !ok || c.N < 1 || c.N > limit {
		return nil, fmt.Errorf("outside fixed experiment bounds")
	}
	source := ""
	switch c.Kind {
	case "depth":
		source = "x=" + strings.Repeat("(", c.N) + "0" + strings.Repeat(")", c.N) + "\n"
	case "broad", "cancel-native", "cancel-walk":
		source = "x=[" + strings.Repeat("0,", c.N) + "]\n"
	case "token":
		source = "x='" + strings.Repeat("a", c.N) + "'\n"
	case "matches":
		source = strings.Repeat("os.open('unexecuted')\n", c.N)
	case "files", "invalid-last":
		source = "os.open('unexecuted')\n"
	}
	count := 1
	if c.Kind == "files" {
		count = c.N
	}
	if c.Kind == "invalid-last" {
		count = 2
	}
	// Conservative aggregate bound includes the invalid trailing file.
	if len(source)*count > 16<<20 {
		return nil, fmt.Errorf("input byte ceiling")
	}
	inputs := make([]probeInput, count)
	for i := range inputs {
		inputs[i] = probeInput{fmt.Sprintf("fixture-%d.py", i), []byte(source)}
	}
	if c.Kind == "invalid-last" {
		inputs[1].Bytes = []byte("return\n")
	}
	return inputs, nil
}

type cappedOutput struct {
	buffer    bytes.Buffer // private: do not promote ReaderFrom and bypass Write
	limit     int
	truncated bool
}

func (w *cappedOutput) Len() int       { return w.buffer.Len() }
func (w *cappedOutput) String() string { return w.buffer.String() }

func (w *cappedOutput) Write(p []byte) (int, error) {
	n := len(p)
	left := w.limit - w.Len()
	if n > left {
		w.truncated = true
	}
	if left > 0 {
		if left < len(p) {
			p = p[:left]
		}
		_, _ = w.buffer.Write(p)
	}
	return n, nil
}

type resourceMeasurement struct {
	Case                        resourceCase
	Bytes, Files, Matches       int
	ElapsedMS, CancelToReturnMS float64
	Outcome, RSS                string
	VmHWMKiB                    int64
}

// Stop=false only means the callback may have started. Its timestamp must
// precede the observed return before treating it as an in-flight cancellation.
func cancellationOutcome(at, returned time.Time, value python.Value, err error) (string, float64, error) {
	if !at.IsZero() && !at.After(returned) {
		if err == nil || value != nil {
			return "", 0, fmt.Errorf("cancelled operation returned candidate")
		}
		return "context cancellation", float64(returned.Sub(at).Microseconds()) / 1000, nil
	}
	if err != nil {
		return "rejection", 0, nil
	}
	return "success before cancellation", 0, nil
}

func validateResourceMeasurement(c resourceCase, m resourceMeasurement) error {
	inputs, err := resourceInputs(c)
	if err != nil {
		return err
	}
	total := 0
	for _, in := range inputs {
		total += len(in.Bytes)
	}
	finite := func(n float64) bool { return !math.IsNaN(n) && !math.IsInf(n, 0) }
	if m.Case != c || m.Bytes != total || m.Files != len(inputs) || m.Matches < 0 || m.Matches > 100000 ||
		!finite(m.ElapsedMS) || m.ElapsedMS <= 0 || m.ElapsedMS > float64(childDeadline.Milliseconds()) ||
		!finite(m.CancelToReturnMS) || m.CancelToReturnMS < 0 || m.CancelToReturnMS > m.ElapsedMS ||
		m.VmHWMKiB < 0 || m.VmHWMKiB > 8<<20 {
		return fmt.Errorf("invalid receipt metrics")
	}
	if (m.RSS != "RSSUnavailable" || m.VmHWMKiB != 0) && (m.RSS != "processRSS VmHWM" || m.VmHWMKiB <= 0) {
		return fmt.Errorf("invalid RSS receipt")
	}
	isCancel := strings.HasPrefix(c.Kind, "cancel-")
	wantMatches := 0
	switch m.Outcome {
	case "success":
		if isCancel || c.Kind == "invalid-last" {
			return fmt.Errorf("inconsistent success")
		}
		if c.Kind == "matches" || c.Kind == "files" {
			wantMatches = c.N
		}
	case "success before cancellation", "context cancellation":
		if !isCancel {
			return fmt.Errorf("unexpected cancellation outcome")
		}
	case "rejection":
	default:
		return fmt.Errorf("unknown child outcome")
	}
	if m.Matches != wantMatches || (m.Outcome != "context cancellation" && m.CancelToReturnMS != 0) {
		return fmt.Errorf("inconsistent result counts or cancellation latency")
	}
	return nil
}

func resourceReceipt(c resourceCase, text string, truncated bool) (resourceMeasurement, error) {
	var result resourceMeasurement
	if truncated || len(text) > 64<<10 {
		return result, fmt.Errorf("truncated or oversized child output")
	}
	count := 0
	for _, line := range strings.SplitAfter(text, "\n") {
		if strings.HasPrefix(line, "RESOURCE ") {
			count++
			if !strings.HasSuffix(line, "\n") {
				return result, fmt.Errorf("incomplete receipt line")
			}
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "RESOURCE ")), &result); err != nil {
				return result, err
			}
			if err := validateResourceMeasurement(c, result); err != nil {
				return result, err
			}
		}
	}
	if count != 1 {
		return result, fmt.Errorf("expected exactly one receipt, got %d", count)
	}
	return result, nil
}

func processRSS() (int64, string) {
	f, err := os.Open("/proc/self/status")
	if err != nil {
		return 0, "RSSUnavailable"
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) == 3 && fields[0] == "VmHWM:" {
			n, err := strconv.ParseInt(fields[1], 10, 64)
			if err == nil && fields[2] == "kB" {
				return n, "processRSS VmHWM"
			}
		}
	}
	return 0, "RSSUnavailable"
}

func TestResourceBoundsPinnedToCeilingsNeverEscalated(t *testing.T) {
	for _, c := range []resourceCase{{"matches", 100001}, {"token", 16 << 20}, {"unknown", 1}, {"files", 1001}, {"depth", 1001}, {"broad", 100001}, {"matches", -1}} {
		if _, err := resourceInputs(c); err == nil {
			t.Fatalf("accepted out-of-contract case: %+v", c)
		}
	}
	for _, c := range resourceCases() {
		inputs, err := resourceInputs(c)
		if err != nil {
			t.Fatal(err)
		}
		total := 0
		for _, in := range inputs {
			total += len(in.Bytes)
		}
		if total > 16<<20 || len(inputs) > 1001 {
			t.Fatal("generator exceeded bounds")
		}
	}
}

// Fixed finite child output exercises os/exec's real io.Copy method selection.
func TestResourceOutputFixture(t *testing.T) {
	if os.Getenv("GCE_RESOURCE_OUTPUT_FIXTURE") != "1" {
		t.Skip("output fixture child only")
	}
	mode := os.Getenv("GCE_RESOURCE_RECEIPT_FIXTURE")
	if mode == "duplicate" || mode == "metrics" {
		m := resourceMeasurement{Case: resourceCase{"matches", 10000}, Bytes: 220000, Files: 1, Matches: 10000, ElapsedMS: 1, Outcome: "success", RSS: "processRSS VmHWM", VmHWMKiB: 1}
		if mode == "metrics" {
			m.Bytes, m.Files, m.Matches, m.ElapsedMS, m.VmHWMKiB, m.Outcome = -7, -9, -1, -1, -10, "fabricated"
		}
		b, _ := json.Marshal(m)
		fmt.Println("RESOURCE " + string(b))
		if mode == "duplicate" {
			fmt.Println("RESOURCE " + string(b))
		}
		return
	}
	fmt.Fprint(os.Stdout, strings.Repeat("x", 100<<10))
	fmt.Fprint(os.Stderr, strings.Repeat("y", 100<<10))
}

func TestResourceOutputCapPinnedToWriteNeverPromotedReadFrom(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe, "-test.run=^TestResourceOutputFixture$")
	cmd.Env = append(os.Environ(), "GCE_RESOURCE_OUTPUT_FIXTURE=1")
	w := &cappedOutput{limit: 64 << 10}
	cmd.Stdout, cmd.Stderr = w, w
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	if _, promoted := any(w).(io.ReaderFrom); promoted {
		t.Error("capped sink unexpectedly exposes ReaderFrom")
	}
	if w.Len() != 64<<10 || !w.truncated {
		t.Fatalf("real subprocess retained %d bytes, truncated=%v; want 65536 and true", w.Len(), w.truncated)
	}
}

func TestResourceReceiptSubprocessNegatives(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"duplicate", "metrics"} {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		cmd := exec.CommandContext(ctx, exe, "-test.run=^TestResourceOutputFixture$")
		cmd.Env = append(os.Environ(), "GCE_RESOURCE_OUTPUT_FIXTURE=1", "GCE_RESOURCE_RECEIPT_FIXTURE="+mode)
		w := &cappedOutput{limit: 64 << 10}
		cmd.Stdout, cmd.Stderr = w, w
		err := cmd.Run()
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := resourceReceipt(resourceCase{"matches", 10000}, w.String(), w.truncated); err == nil {
			t.Fatalf("accepted subprocess %s receipts", mode)
		}
	}
}

func TestResourceCancellationAfterReturnNative(t *testing.T) {
	if os.Getenv("GCE_RESOURCE_AFTER_RETURN_FIXTURE") != "1" {
		exe, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), childDeadline)
		defer cancel()
		cmd := exec.CommandContext(ctx, exe, "-test.run=^TestResourceCancellationAfterReturnNative$", "-test.v")
		cmd.Env = append(os.Environ(), "GCE_RESOURCE_AFTER_RETURN_FIXTURE=1")
		w := &cappedOutput{limit: 64 << 10}
		cmd.Stdout, cmd.Stderr = w, w
		if err := cmd.Run(); err != nil || w.truncated {
			t.Fatalf("native fixture: %v %s", err, w.String())
		}
		return
	}
	p, err := python.New(python.Config{})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	ast, err := p.Import(context.Background(), "ast")
	if err != nil {
		t.Fatal(err)
	}
	inputs, err := resourceInputs(resourceCase{"cancel-native", 100000})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stamps := make(chan time.Time, 1)
	timer := time.AfterFunc(time.Hour, func() { at := time.Now(); cancel(); stamps <- at })
	defer timer.Stop()
	value, err := ast.CallMethod(ctx, "parse", python.ValueOf(inputs[0].Bytes))
	returned := time.Now()
	timer.Reset(0) // controlled release AFTER the actual native parse returns
	at := <-stamps
	if err != nil || value == nil || !at.After(returned) || timer.Stop() {
		t.Fatal("native after-return premise failed")
	}
	outcome, latency, e := cancellationOutcome(at, returned, value, err)
	if e != nil || outcome != "success before cancellation" || latency != 0 {
		t.Fatalf("controlled native completion: %s %f %v", outcome, latency, e)
	}
	t.Log("controlled after-return callback; no native candidate leak")
}

func TestResourceCancellationOrdering(t *testing.T) {
	returned := time.Now()
	for _, after := range []bool{false, true} {
		at := returned.Add(-time.Millisecond)
		var value python.Value
		err := error(context.Canceled)
		if after {
			at = returned.Add(time.Millisecond)
			err, value = nil, python.ValueOf(1)
		}
		outcome, latency, e := cancellationOutcome(at, returned, value, err)
		want := "context cancellation"
		if after {
			want = "success before cancellation"
		}
		if e != nil || outcome != want || latency < 0 || (after && latency != 0) {
			t.Fatalf("after=%v: %s %f %v", after, outcome, latency, e)
		}
	}
	// Real callback has started, but cancellation is held until AFTER the
	// observed return. Stop=false is not an ordering proof.
	release := make(chan struct{})
	started := make(chan struct{})
	stamps := make(chan time.Time, 1)
	timer := time.AfterFunc(0, func() { close(started); <-release; stamps <- time.Now() })
	<-started
	returned = time.Now()
	close(release)
	at := <-stamps
	if timer.Stop() || !at.After(returned) {
		t.Fatal("controlled after-return premise failed")
	}
	if outcome, ms, err := cancellationOutcome(at, returned, python.ValueOf(1), nil); err != nil || outcome != "success before cancellation" || ms != 0 {
		t.Fatalf("after-return callback misclassified: %s %f %v", outcome, ms, err)
	}
	if _, _, err := cancellationOutcome(returned.Add(-time.Millisecond), returned, python.ValueOf(1), nil); err == nil {
		t.Fatal("pre-return cancellation accepted a candidate")
	}
}

func TestResourceReceiptValidationPinnedToOneNeverDuplicate(t *testing.T) {
	c := resourceCase{"matches", 10000}
	good := resourceMeasurement{Case: c, Bytes: 220000, Files: 1, Matches: 10000, ElapsedMS: 1, Outcome: "success", RSS: "processRSS VmHWM", VmHWMKiB: 1}
	encode := func(m resourceMeasurement) string {
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		return "RESOURCE " + string(b) + "\n"
	}
	line := encode(good)
	for _, text := range []string{line, "=== RUN child\n" + line + "PASS\n"} {
		if _, err := resourceReceipt(c, text, false); err != nil {
			t.Fatal(err)
		}
	}
	bad := good
	bad.Bytes, bad.Files, bad.Matches, bad.ElapsedMS, bad.VmHWMKiB, bad.Outcome = -7, -9, -1, -1, -10, "fabricated"
	for _, text := range []string{"", line + line, encode(bad), "RESOURCE {\n", strings.TrimSuffix(line, "\n"), strings.Repeat("x", 64<<10) + line} {
		if _, err := resourceReceipt(c, text, false); err == nil {
			t.Fatalf("accepted malformed receipt: %.80s", text)
		}
	}
	if _, err := resourceReceipt(c, line, true); err == nil {
		t.Fatal("accepted truncated output")
	}
	for _, mutate := range []func(*resourceMeasurement){
		func(m *resourceMeasurement) { m.Bytes++ }, func(m *resourceMeasurement) { m.Files++ },
		func(m *resourceMeasurement) { m.Matches-- }, func(m *resourceMeasurement) { m.Outcome = "fabricated" },
		func(m *resourceMeasurement) { m.ElapsedMS = -1 }, func(m *resourceMeasurement) { m.CancelToReturnMS = -1 },
		func(m *resourceMeasurement) { m.VmHWMKiB = -1 }, func(m *resourceMeasurement) { m.RSS = "fabricated" },
		func(m *resourceMeasurement) { m.ElapsedMS = math.NaN() }, func(m *resourceMeasurement) { m.CancelToReturnMS = math.Inf(1) },
		func(m *resourceMeasurement) { m.ElapsedMS = math.Inf(-1) },
	} {
		m := good
		mutate(&m)
		if err := validateResourceMeasurement(c, m); err == nil {
			t.Fatalf("accepted invalid metrics: %+v", m)
		}
	}
}

func TestResourceHarnessContract(t *testing.T) {
	if childDeadline.Seconds() != 30 || len(resourceCases()) != 13 {
		t.Fatal("unbounded harness contract")
	}
	w := &cappedOutput{limit: 8}
	n, err := w.Write([]byte("0123456789"))
	if n != 10 || err != nil || w.String() != "01234567" {
		t.Fatal("output not capped")
	}
	for _, c := range []resourceCase{{"depth", 1}, {"broad", 1}, {"token", 1}, {"matches", 1}, {"files", 1}} {
		inputs, err := resourceInputs(c)
		want := map[string]string{"depth": "x=(0)\n", "broad": "x=[0,]\n", "token": "x='a'\n", "matches": "os.open('unexecuted')\n", "files": "os.open('unexecuted')\n"}
		if err != nil || len(inputs) != 1 || string(inputs[0].Bytes) != want[c.Kind] {
			t.Fatalf("generator shape drift: %+v", c)
		}
	}
	inputs, err := resourceInputs(resourceCase{"invalid-last", 1})
	if err != nil || len(inputs) != 2 || string(inputs[1].Bytes) != "return\n" {
		t.Fatal("invalid-last fixture drift")
	}
}
