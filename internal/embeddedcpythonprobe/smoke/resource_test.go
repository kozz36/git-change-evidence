package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	python "github.com/goccy/go-python"
)

// Explicit test selection OR GCE_RESOURCE_STRESS=1 opts in; default go test
// skips stress. The child selector is never accepted as a generated parameter.
func TestResourceStress(t *testing.T) {
	if os.Getenv("GCE_RESOURCE_STRESS") != "1" && !strings.Contains(strings.Join(os.Args, " "), "TestResourceStress") {
		t.Skip("opt-in resource experiment")
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range resourceCases() { // serial fresh processes, at most 13*30s
		ctx, cancel := context.WithTimeout(context.Background(), childDeadline)
		cmd := exec.CommandContext(ctx, exe, "-test.run=^TestResourceChild$", "-test.v")
		cmd.Env = append(os.Environ(), fmt.Sprintf("GCE_RESOURCE_CHILD=%d", i))
		output := &cappedOutput{limit: 64 << 10}
		cmd.Stdout = output
		cmd.Stderr = output
		started := time.Now()
		err := cmd.Run()
		killed := ctx.Err() != nil
		cancel()
		if killed {
			inputs, _ := resourceInputs(c)
			m := resourceMeasurement{Case: c, Files: len(inputs), Outcome: "timeout-killed", RSS: "RSSUnavailable", ElapsedMS: float64(time.Since(started).Microseconds()) / 1000}
			for _, in := range inputs {
				m.Bytes += len(in.Bytes)
			}
			b, _ := json.Marshal(m)
			t.Log("RESOURCE " + string(b))
			continue
		}
		if err != nil {
			t.Fatalf("child %+v: %v: %s", c, err, output.String())
		}
		m, err := resourceReceipt(c, output.String(), output.truncated)
		if err != nil {
			t.Fatalf("invalid child receipt: %v", err)
		}
		b, _ := json.Marshal(m)
		t.Log("RESOURCE " + string(b))
	}
}

func TestResourceChild(t *testing.T) {
	selector := os.Getenv("GCE_RESOURCE_CHILD")
	if selector == "" {
		t.Skip("child only")
	}
	i, err := strconv.Atoi(selector)
	cases := resourceCases()
	if err != nil || i < 0 || i >= len(cases) {
		t.Fatal("invalid child selector")
	}
	c := cases[i]
	inputs, err := resourceInputs(c)
	if err != nil {
		t.Fatal(err)
	}
	m := resourceMeasurement{Case: c, Files: len(inputs), Outcome: "success"}
	for _, in := range inputs {
		m.Bytes += len(in.Bytes)
	}
	started := time.Now()
	p, err := python.New(python.Config{})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if strings.HasPrefix(c.Kind, "cancel-") {
		// Initialize/import/prepare tree before cancellation, so native parsing and
		// Python-level walking are distinct operations. No scanned code is run.
		ast, err := p.Import(context.Background(), "ast")
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var tree python.Value
		if c.Kind == "cancel-walk" {
			tree, err = ast.CallMethod(context.Background(), "parse", python.ValueOf(inputs[0].Bytes))
			if err != nil {
				t.Fatal(err)
			}
			r, e := p.Eval(context.Background(), exactCallsHelper)
			if e != nil || r.Error != nil {
				t.Fatalf("helper: %v %v", e, r.Error)
			}
		}
		main, err := p.Import(context.Background(), "__main__")
		if err != nil {
			t.Fatal(err)
		}
		cancelled := make(chan time.Time, 1)
		timer := time.AfterFunc(cancelDelay, func() { at := time.Now(); cancel(); cancelled <- at })
		defer timer.Stop()
		var value python.Value
		if c.Kind == "cancel-native" {
			value, err = ast.CallMethod(ctx, "parse", python.ValueOf(inputs[0].Bytes))
		} else {
			value, err = main.CallMethod(ctx, "_count_exact_calls", tree)
		}
		returned := time.Now()
		var at time.Time
		if !timer.Stop() {
			at = <-cancelled // wait for callback completion, then compare timestamps
		}
		m.Outcome, m.CancelToReturnMS, err = cancellationOutcome(at, returned, value, err)
		if err != nil {
			t.Fatal(err)
		}
		// Independently verify terminal batch cancellation discards earlier results.
		cancel()
		batch, e := probeBatch(ctx, p, inputs)
		if e == nil || batch != nil {
			t.Fatal("cancelled batch returned partial results")
		}
	} else {
		files, e := probeBatch(context.Background(), p, inputs)
		if e != nil {
			m.Outcome = "rejection"
			if files != nil {
				t.Fatal("rejection returned partial batch")
			}
		} else {
			for _, f := range files {
				m.Matches += len(f.Calls)
			}
		}
		if c.Kind == "invalid-last" && (e == nil || files != nil) {
			t.Fatal("invalid-last batch accepted")
		}
		if c.Kind == "matches" && e == nil && m.Matches != c.N {
			t.Fatal("match count drift")
		}
	}
	m.ElapsedMS = float64(time.Since(started).Microseconds()) / 1000
	m.VmHWMKiB, m.RSS = processRSS()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("RESOURCE " + string(b))
}
