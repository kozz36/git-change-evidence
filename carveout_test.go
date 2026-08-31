package changeevidence

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strings"
	"testing"
)

func TestPrepareCarveoutReportsInvalidBounds(t *testing.T) {
	for _, test := range []struct {
		name   string
		bounds CarveoutBounds
		limit  CarveoutLimit
	}{
		{"zero byte cap", CarveoutBounds{LineCap: 1, MaxLineBytes: 1}, CarveoutByteLimit},
		{"negative byte cap", CarveoutBounds{ByteCap: -1, LineCap: 1, MaxLineBytes: 1}, CarveoutByteLimit},
		{"unsafe byte cap", CarveoutBounds{ByteCap: CarveoutSafeMaxBytes + 1, LineCap: 1, MaxLineBytes: 1}, CarveoutByteLimit},
		{"zero line cap", CarveoutBounds{ByteCap: 1, MaxLineBytes: 1}, CarveoutLineLimit},
		{"negative line cap", CarveoutBounds{ByteCap: 1, LineCap: -1, MaxLineBytes: 1}, CarveoutLineLimit},
		{"unsafe line cap", CarveoutBounds{ByteCap: 1, LineCap: CarveoutSafeMaxLines + 1, MaxLineBytes: 1}, CarveoutLineLimit},
		{"zero line bytes", CarveoutBounds{ByteCap: 1, LineCap: 1}, CarveoutLineByteLimit},
		{"negative line bytes", CarveoutBounds{ByteCap: 1, LineCap: 1, MaxLineBytes: -1}, CarveoutLineByteLimit},
		{"unsafe line bytes", CarveoutBounds{ByteCap: 1, LineCap: 1, MaxLineBytes: CarveoutSafeMaxLineBytes + 1}, CarveoutLineByteLimit},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertPreparation(t, PrepareCarveout(availableCarveout("x", "y"), test.bounds), CarveoutInvalidBound, test.limit)
		})
	}
}

func TestPrepareCarveoutDistinguishesTypedStates(t *testing.T) {
	available := availableCarveout("x", "")
	unavailableBefore, unavailableAfter := available, available
	unavailableBefore.Before.Available, unavailableAfter.After.Available = false, false
	for _, test := range []struct {
		name    string
		input   CarveoutInput
		bounds  CarveoutBounds
		state   CarveoutState
		limit   CarveoutLimit
		content CarveoutContentName
	}{
		{"prepared", available, validCarveoutBounds(), CarveoutPrepared, "", ""},
		{"invalid bound", unavailableBefore, CarveoutBounds{LineCap: 1, MaxLineBytes: 1}, CarveoutInvalidBound, CarveoutByteLimit, ""},
		{"unavailable before", unavailableBefore, validCarveoutBounds(), CarveoutUnavailableRequiredContent, "", CarveoutBeforeContent},
		{"unavailable after", unavailableAfter, validCarveoutBounds(), CarveoutUnavailableRequiredContent, "", CarveoutAfterContent},
		{"resource bound", availableCarveout("x", "y"), CarveoutBounds{ByteCap: 1, LineCap: 2, MaxLineBytes: 1}, CarveoutResourceBoundExceeded, CarveoutByteLimit, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := PrepareCarveout(test.input, test.bounds)
			assertPreparation(t, got, test.state, test.limit)
			if got.Content != test.content {
				t.Fatalf("content = %q, want %q", got.Content, test.content)
			}
		})
	}
}

func TestPrepareCarveoutHonorsLimitBoundaries(t *testing.T) {
	for _, test := range []struct {
		name, before, after string
		bounds              CarveoutBounds
		state               CarveoutState
		limit               CarveoutLimit
	}{
		{"byte cap exact", "a", "bc", CarveoutBounds{ByteCap: 3, LineCap: 2, MaxLineBytes: 2}, CarveoutPrepared, ""},
		{"byte cap one over", "a", "bc", CarveoutBounds{ByteCap: 2, LineCap: 2, MaxLineBytes: 2}, CarveoutResourceBoundExceeded, CarveoutByteLimit},
		{"line cap exact", "a\nb", "c", CarveoutBounds{ByteCap: 5, LineCap: 3, MaxLineBytes: 1}, CarveoutPrepared, ""},
		{"line cap one over", "a\nb", "c", CarveoutBounds{ByteCap: 5, LineCap: 2, MaxLineBytes: 1}, CarveoutResourceBoundExceeded, CarveoutLineLimit},
		{"line bytes exact", "abc", "", CarveoutBounds{ByteCap: 3, LineCap: 1, MaxLineBytes: 3}, CarveoutPrepared, ""},
		{"line bytes one over", "abc", "", CarveoutBounds{ByteCap: 3, LineCap: 1, MaxLineBytes: 2}, CarveoutResourceBoundExceeded, CarveoutLineByteLimit},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertPreparation(t, PrepareCarveout(availableCarveout(test.before, test.after), test.bounds), test.state, test.limit)
		})
	}
}

func TestPrepareCarveoutNormalizesBlankLines(t *testing.T) {
	got := PrepareCarveout(availableCarveout("\n\nfirst\n\nsecond\n", "\nthird\n\n"), validCarveoutBounds())
	assertPreparation(t, got, CarveoutPrepared, "")
	assertLines(t, got.Prepared, []string{"first", "second"}, []string{"third"})
}

func TestPreparedCarveoutPreservesRawIdentity(t *testing.T) {
	identities := map[Digest]bool{}
	for _, raw := range []string{"", "\n", "\n\n", "\r\n", "\r\n\r\n"} {
		got := PrepareCarveout(availableCarveout(raw, raw), validCarveoutBounds())
		assertPreparation(t, got, CarveoutPrepared, "")
		if got.Prepared.BeforeLineCount() != 0 || got.Prepared.AfterLineCount() != 0 || got.Prepared.BeforeRawDigest() != carveoutDigest(raw) || got.Prepared.AfterRawDigest() != carveoutDigest(raw) {
			t.Fatalf("raw %q preserved nonblank lines or wrong digest", raw)
		}
		if identities[got.Prepared.BeforeRawDigest()] {
			t.Fatalf("duplicate raw digest for %q", raw)
		}
		identities[got.Prepared.BeforeRawDigest()] = true
	}
}

func TestPreparedCarveoutPreservesInvalidUTF8AndUnicodeBytes(t *testing.T) {
	raw := string([]byte{0xff, '\n', 0xe2, 0x98, 0x83})
	got := PrepareCarveout(availableCarveout(raw, ""), validCarveoutBounds())
	assertPreparation(t, got, CarveoutPrepared, "")
	assertLines(t, got.Prepared, []string{string([]byte{0xff}), "☃"}, nil)
	if got.Prepared.BeforeRawDigest() != carveoutDigest(raw) {
		t.Fatalf("invalid UTF-8 digest = %q", got.Prepared.BeforeRawDigest())
	}
}

func TestPrepareCarveoutUsesImmutableStringInput(t *testing.T) {
	raw := "first\nsecond"
	input := availableCarveout(raw, "")
	result := make(chan CarveoutPreparation, 1)
	go func() { result <- PrepareCarveout(input, validCarveoutBounds()) }()
	raw = "changed"
	got := <-result
	assertPreparation(t, got, CarveoutPrepared, "")
	assertLines(t, got.Prepared, []string{"first", "second"}, nil)
	if raw != "changed" {
		t.Fatal("string reassignment failed")
	}
}

var carveoutLineSink string
var carveoutOKSink bool

func TestPreparedCarveoutLineAccessIsSafeAndAllocationFree(t *testing.T) {
	prepared := PrepareCarveout(availableCarveout("first\nsecond", "third"), validCarveoutBounds()).Prepared
	for _, test := range []struct {
		index int
		line  string
		ok    bool
	}{{-1, "", false}, {0, "first", true}, {1, "second", true}, {2, "", false}} {
		if line, ok := prepared.BeforeLine(test.index); line != test.line || ok != test.ok {
			t.Fatalf("BeforeLine(%d) = %q, %v", test.index, line, ok)
		}
	}
	line, _ := prepared.BeforeLine(0)
	line = "changed"
	if again, _ := prepared.BeforeLine(0); line == again {
		t.Fatalf("mutable descriptor view = %q", again)
	}
	for _, index := range []int{-1, 1} {
		if line, ok := prepared.AfterLine(index); line != "" || ok {
			t.Fatalf("AfterLine(%d) = %q, %v", index, line, ok)
		}
	}
	if allocations := testing.AllocsPerRun(100, func() {
		for index := -1; index <= prepared.BeforeLineCount(); index++ {
			carveoutLineSink, carveoutOKSink = prepared.BeforeLine(index)
		}
		for index := -1; index <= prepared.AfterLineCount(); index++ {
			carveoutLineSink, carveoutOKSink = prepared.AfterLine(index)
		}
	}); allocations != 0 {
		t.Fatalf("line access allocations = %v, want 0", allocations)
	}
}

func TestPrepareCarveoutReportsHostileInputResourceBoundsBeforeAllocation(t *testing.T) {
	newlineHeavy := availableCarveout(strings.Repeat("x\n", 4096), "")
	bounds := CarveoutBounds{ByteCap: len(newlineHeavy.Before.Raw), LineCap: 2, MaxLineBytes: 1}
	if allocations := testing.AllocsPerRun(10, func() {
		assertPreparation(t, PrepareCarveout(newlineHeavy, bounds), CarveoutResourceBoundExceeded, CarveoutLineLimit)
	}); allocations != 0 {
		t.Fatalf("newline-heavy rejection allocations = %v, want 0", allocations)
	}
	veryLong := availableCarveout(strings.Repeat("x", 4096), "")
	assertPreparation(t, PrepareCarveout(veryLong, CarveoutBounds{ByteCap: 4096, LineCap: 1, MaxLineBytes: 4095}), CarveoutResourceBoundExceeded, CarveoutLineByteLimit)
}

func TestWithinCarveoutByteCapIsOverflowSafe(t *testing.T) {
	for _, test := range []struct {
		left, right, limit int
		want               bool
	}{{1, 2, 3, true}, {1, 3, 3, false}, {math.MaxInt - 1, 1, math.MaxInt, true}, {math.MaxInt, 1, math.MaxInt, false}} {
		if got := withinCarveoutByteCap(test.left, test.right, test.limit); got != test.want {
			t.Fatalf("withinCarveoutByteCap(%d, %d, %d) = %v, want %v", test.left, test.right, test.limit, got, test.want)
		}
	}
}

func availableCarveout(before, after string) CarveoutInput {
	return CarveoutInput{Before: CarveoutContent{Available: true, Raw: before}, After: CarveoutContent{Available: true, Raw: after}}
}
func validCarveoutBounds() CarveoutBounds {
	return CarveoutBounds{ByteCap: 1024, LineCap: 128, MaxLineBytes: 128}
}
func carveoutDigest(raw string) Digest {
	sum := sha256.Sum256([]byte(raw))
	return Digest(hex.EncodeToString(sum[:]))
}

func assertLines(t *testing.T, got PreparedCarveout, before, after []string) {
	t.Helper()
	if got.BeforeLineCount() != len(before) || got.AfterLineCount() != len(after) {
		t.Fatalf("line counts = %d, %d; want %d, %d", got.BeforeLineCount(), got.AfterLineCount(), len(before), len(after))
	}
	actualBefore, actualAfter := make([]string, len(before)), make([]string, len(after))
	for index := range actualBefore {
		actualBefore[index], _ = got.BeforeLine(index)
	}
	for index := range actualAfter {
		actualAfter[index], _ = got.AfterLine(index)
	}
	for index, want := range before {
		if actualBefore[index] != want {
			t.Fatalf("before line %d = %q, want %q", index, actualBefore[index], want)
		}
	}
	for index, want := range after {
		if actualAfter[index] != want {
			t.Fatalf("after line %d = %q, want %q", index, actualAfter[index], want)
		}
	}
}

func assertPreparation(t *testing.T, got CarveoutPreparation, state CarveoutState, limit CarveoutLimit) {
	t.Helper()
	if got.State != state || got.Limit != limit || (state != CarveoutPrepared && (got.Prepared.BeforeLineCount() != 0 || got.Prepared.AfterLineCount() != 0)) {
		t.Fatalf("preparation = %#v, want state %q limit %q with no partial data on failure", got, state, limit)
	}
}
