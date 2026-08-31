package changeevidence_test

import (
	changeevidence "github.com/kozz36/git-change-evidence"
	"reflect"
	"strings"
	"testing"
	"time"
)

type matcherClock struct {
	times []time.Time
	last  time.Time
}

func (c *matcherClock) Now() time.Time {
	if len(c.times) != 0 {
		c.last, c.times = c.times[0], c.times[1:]
	}
	return c.last
}
func deadlineClock(at int) *matcherClock {
	deadline := time.Unix(0, 0).Add(time.Millisecond)
	clock := &matcherClock{times: make([]time.Time, at), last: deadline}
	clock.times[at-1] = deadline
	return clock
}
func prepared(before, after string) changeevidence.PreparedCarveout {
	return changeevidence.PrepareCarveout(changeevidence.CarveoutInput{Before: changeevidence.CarveoutContent{Available: true, Raw: before}, After: changeevidence.CarveoutContent{Available: true, Raw: after}}, changeevidence.CarveoutBounds{ByteCap: 1024, LineCap: 128, MaxLineBytes: 128}).Prepared
}
func matcherOptions(clock *matcherClock) changeevidence.CarveoutMatcherOptions {
	return changeevidence.CarveoutMatcherOptions{PairCap: 32, TimeCap: time.Hour, Clock: clock}
}
func TestMatchPreparedCarveoutMatchesDeterministically(t *testing.T) {
	input := prepared("same\nsame\nabcde\nabXde", "same\nabYde\nsame\nabcde")
	first, second := changeevidence.MatchPreparedCarveout(input, matcherOptions(&matcherClock{})), changeevidence.MatchPreparedCarveout(input, matcherOptions(&matcherClock{}))
	want := []changeevidence.CarveoutLineMatch{{BeforeIndex: 0, AfterIndex: 0}, {BeforeIndex: 3, AfterIndex: 1}, {BeforeIndex: 1, AfterIndex: 2}, {BeforeIndex: 2, AfterIndex: 3}}
	for _, got := range []changeevidence.CarveoutMatchAnalysis{first, second} {
		if got.State != changeevidence.CarveoutMatchCompleted || got.Result == nil || got.Result.Moved() != 4 || got.Result.MatchCount() != 4 {
			t.Fatalf("analysis = %#v", got)
		}
		for index, want := range want {
			if match, ok := got.Result.MatchAt(index); !ok || match != want {
				t.Fatalf("MatchAt(%d) = %#v, %v; want %#v", index, match, ok, want)
			}
		}
	}
}
func TestMatchPreparedCarveoutDeadlineBoundaries(t *testing.T) {
	long := strings.Repeat("x", 96)
	for _, test := range []struct {
		name, before, after string
		at, inspections     int
	}{
		{"before allocation", "x", "x", 2, 0}, {"after exact comparison", long + "a", long + "b", 8, 1},
		{"fuzzy row allocation", "a", "b", 13, 2}, {"fuzzy row initialization", "a", "b", 14, 2}, {"fuzzy cell", "a", "b", 17, 2},
		{"result assembly", "x", "x", 10, 1}, {"before completion", "x", "x", 11, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			o := matcherOptions(deadlineClock(test.at))
			o.TimeCap = time.Millisecond
			got := changeevidence.MatchPreparedCarveout(prepared(test.before, test.after), o)
			if got.State != changeevidence.CarveoutMatchTimeout || got.Provenance.State != changeevidence.CarveoutMatchTimeout || got.Provenance.Reason != changeevidence.CarveoutMatchDeadlineExceeded || got.Provenance.PairInspections != test.inspections || got.Result != nil || got.Fallback != nil {
				t.Fatalf("analysis = %#v", got)
			}
		})
	}
}
func TestMatchPreparedCarveoutRequiresMatches(t *testing.T) {
	valid := &changeevidence.CarveoutMatchFallback{Audited: true, Reason: "audit", Evidence: "digest"}
	for _, test := range []struct {
		name     string
		required bool
		cap      int
		fallback *changeevidence.CarveoutMatchFallback
		state    changeevidence.CarveoutMatchState
		reason   changeevidence.CarveoutMatchReason
		retained bool
	}{
		{"optional", false, 32, nil, changeevidence.CarveoutMatchCompleted, "", false}, {"completed ignores fallback", false, 32, valid, changeevidence.CarveoutMatchCompleted, "", false}, {"required", true, 32, nil, changeevidence.CarveoutMatchImpossibility, changeevidence.CarveoutMatchRequiredMatchUnavailable, false},
		{"audited fallback", true, 32, valid, changeevidence.CarveoutMatchFallbackApplied, changeevidence.CarveoutMatchRequiredMatchUnavailable, true}, {"unaudited fallback", true, 32, &changeevidence.CarveoutMatchFallback{Reason: "audit", Evidence: "digest"}, changeevidence.CarveoutMatchImpossibility, changeevidence.CarveoutMatchRequiredMatchUnavailable, false}, {"invalid fallback", true, 32, &changeevidence.CarveoutMatchFallback{Audited: true, Reason: "audit"}, changeevidence.CarveoutMatchImpossibility, changeevidence.CarveoutMatchRequiredMatchUnavailable, false},
		{"negative fallback moved", true, 32, &changeevidence.CarveoutMatchFallback{Audited: true, Moved: -1, Reason: "audit", Evidence: "digest"}, changeevidence.CarveoutMatchImpossibility, changeevidence.CarveoutMatchRequiredMatchUnavailable, false}, {"oversized fallback moved", true, 32, &changeevidence.CarveoutMatchFallback{Audited: true, Moved: 2, Reason: "audit", Evidence: "digest"}, changeevidence.CarveoutMatchImpossibility, changeevidence.CarveoutMatchRequiredMatchUnavailable, false}, {"invalid options", true, 0, valid, changeevidence.CarveoutMatchFallbackApplied, changeevidence.CarveoutMatchInvalidPairCap, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			o := matcherOptions(&matcherClock{})
			o.RequireMatch, o.PairCap, o.Fallback = test.required, test.cap, test.fallback
			got := changeevidence.MatchPreparedCarveout(prepared("x", "y"), o)
			if got.State != test.state || got.Provenance.Reason != test.reason || got.Provenance.State != map[bool]changeevidence.CarveoutMatchState{true: changeevidence.CarveoutMatchImpossibility, false: test.state}[test.retained] || (got.Result != nil) != (test.state == changeevidence.CarveoutMatchCompleted) || got.Result.MatchCount() != 0 || got.Result.Moved() != 0 || (got.Fallback != nil) != test.retained {
				t.Fatalf("analysis = %#v", got)
			}
			if test.retained && *got.Fallback != *valid {
				t.Fatalf("fallback = %#v", got.Fallback)
			}
		})
	}
	o := matcherOptions(&matcherClock{})
	o.RequireMatch, o.Fallback = true, valid
	got := changeevidence.MatchPreparedCarveout(prepared("x", "y"), o)
	valid.Reason, valid.Evidence = "changed", "changed"
	if got.Fallback == nil || got.Fallback.Reason != "audit" || got.Fallback.Evidence != "digest" {
		t.Fatalf("retained fallback = %#v", got.Fallback)
	}
}

func TestMatchPreparedCarveoutResultAPIIsImmutable(t *testing.T) {
	result := changeevidence.MatchPreparedCarveout(prepared("x", "x"), matcherOptions(&matcherClock{})).Result
	moved := result.Moved()
	if _, exported := reflect.TypeOf(*result).FieldByName("Moved"); exported || moved != 1 {
		t.Fatalf("result = %#v", result)
	}
	match, _ := result.MatchAt(0)
	match.BeforeIndex = 99
	if again, _ := result.MatchAt(0); again.BeforeIndex != 0 || result.Moved() != moved {
		t.Fatalf("mutable match = %#v", again)
	}
	if allocations := testing.AllocsPerRun(100, func() {
		for index := -1; index <= result.MatchCount(); index++ {
			_, _ = result.MatchAt(index)
		}
	}); allocations != 0 {
		t.Fatalf("MatchAt allocations = %v", allocations)
	}
}

func TestMatchPreparedCarveoutPreservesFuzzyAndCapBounds(t *testing.T) {
	invalidUTF8 := string([]byte{0xff}) + "☃"
	valid := &changeevidence.CarveoutMatchFallback{Audited: true, Moved: 1, Reason: "audit", Evidence: "digest"}
	for _, test := range []struct {
		before, after string
		cap, want     int
		fallback      *changeevidence.CarveoutMatchFallback
		state         changeevidence.CarveoutMatchState
	}{
		{strings.Repeat("a", 100), strings.Repeat("a", 80) + strings.Repeat("b", 20), 32, 1, nil, changeevidence.CarveoutMatchCompleted}, {strings.Repeat("a", 100), strings.Repeat("a", 79) + strings.Repeat("b", 21), 32, 0, nil, changeevidence.CarveoutMatchCompleted}, {invalidUTF8 + "a", invalidUTF8 + "b", 32, 1, nil, changeevidence.CarveoutMatchCompleted},
		{"same\nsame", "same\nsame\nsame", 32, 2, nil, changeevidence.CarveoutMatchCompleted}, {"x\ny", "y", 2, 1, nil, changeevidence.CarveoutMatchCompleted}, {"x\ny", "y", 1, 0, nil, changeevidence.CarveoutMatchTimeout}, {"x\ny", "y", 1, 0, valid, changeevidence.CarveoutMatchFallbackApplied},
	} {
		o := matcherOptions(&matcherClock{})
		o.PairCap, o.Fallback = test.cap, test.fallback
		got := changeevidence.MatchPreparedCarveout(prepared(test.before, test.after), o)
		if got.State != test.state || (got.Result == nil) != (test.state != changeevidence.CarveoutMatchCompleted) || got.Result.MatchCount() != test.want || (got.Fallback != nil) != (test.state == changeevidence.CarveoutMatchFallbackApplied) || test.state == changeevidence.CarveoutMatchFallbackApplied && got.Fallback.Moved != 1 || test.cap == 1 && (got.Result != nil || got.Provenance.State != changeevidence.CarveoutMatchTimeout || got.Provenance.Reason != changeevidence.CarveoutMatchPairCapExceeded || got.Provenance.PairInspections != 1) {
			t.Fatalf("analysis = %#v", got)
		}
	}
}
