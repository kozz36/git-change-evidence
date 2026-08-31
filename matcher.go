package changeevidence

import (
	"strings"
	"time"
)

type CarveoutClock interface{ Now() time.Time }
type CarveoutMatcherOptions struct {
	PairCap      int
	TimeCap      time.Duration
	Clock        CarveoutClock
	RequireMatch bool
	Fallback     *CarveoutMatchFallback
}
type CarveoutMatchState string

const (
	CarveoutMatchCompleted       CarveoutMatchState = "completed"
	CarveoutMatchTimeout         CarveoutMatchState = "timeout"
	CarveoutMatchImpossibility   CarveoutMatchState = "impossibility"
	CarveoutMatchFallbackApplied CarveoutMatchState = "fallback_applied"
)

type CarveoutMatchReason string

const (
	CarveoutMatchDeadlineExceeded         CarveoutMatchReason = "deadline_exceeded"
	CarveoutMatchPairCapExceeded          CarveoutMatchReason = "pair_cap_exceeded"
	CarveoutMatchInvalidPairCap           CarveoutMatchReason = "invalid_pair_cap"
	CarveoutMatchInvalidTimeCap           CarveoutMatchReason = "invalid_time_cap"
	CarveoutMatchMissingClock             CarveoutMatchReason = "missing_clock"
	CarveoutMatchRequiredMatchUnavailable CarveoutMatchReason = "required_match_unavailable"
)

type CarveoutLineMatch struct{ BeforeIndex, AfterIndex int }
type CarveoutMatchResult struct {
	moved   int
	matches []CarveoutLineMatch
}

func (r *CarveoutMatchResult) Moved() int {
	if r == nil {
		return 0
	}
	return r.moved
}
func (r *CarveoutMatchResult) MatchCount() int {
	if r == nil {
		return 0
	}
	return len(r.matches)
}
func (r *CarveoutMatchResult) MatchAt(index int) (CarveoutLineMatch, bool) {
	if r == nil || index < 0 || index >= len(r.matches) {
		return CarveoutLineMatch{}, false
	}
	return r.matches[index], true
}

type CarveoutMatchProvenance struct {
	State           CarveoutMatchState
	Reason          CarveoutMatchReason
	PairInspections int
}
type CarveoutMatchFallback struct {
	Audited          bool
	Moved            int
	Reason, Evidence string
}
type CarveoutMatchAnalysis struct {
	State      CarveoutMatchState
	Result     *CarveoutMatchResult
	Provenance CarveoutMatchProvenance
	Fallback   *CarveoutMatchFallback
}
type preparedCarveoutMatcher struct {
	clock                CarveoutClock
	deadline             time.Time
	pairCap, inspections int
	reason               CarveoutMatchReason
}

func MatchPreparedCarveout(prepared PreparedCarveout, options CarveoutMatcherOptions) CarveoutMatchAnalysis {
	if reason := matchingPrecondition(options); reason != "" {
		return matchingFailure(prepared, CarveoutMatchImpossibility, reason, 0, options.Fallback)
	}
	matcher := preparedCarveoutMatcher{clock: options.Clock, deadline: options.Clock.Now().Add(options.TimeCap), pairCap: options.PairCap}
	if !matcher.check() {
		return matcher.failure(prepared, options.Fallback)
	}
	before, after := prepared.BeforeLineCount(), prepared.AfterLineCount()
	consumed := make([]bool, before)
	if !matcher.check() {
		return matcher.failure(prepared, options.Fallback)
	}
	selected := make([]int, after)
	if !matcher.check() {
		return matcher.failure(prepared, options.Fallback)
	}
	for index := range selected {
		if !matcher.check() {
			return matcher.failure(prepared, options.Fallback)
		}
		selected[index] = -1
	}
	if !matcher.scan(prepared, consumed, selected, false) || !matcher.scan(prepared, consumed, selected, true) {
		return matcher.failure(prepared, options.Fallback)
	}
	if !matcher.check() {
		return matcher.failure(prepared, options.Fallback)
	}
	matches := make([]CarveoutLineMatch, 0, before)
	for afterIndex, beforeIndex := range selected {
		if !matcher.check() {
			return matcher.failure(prepared, options.Fallback)
		}
		if beforeIndex >= 0 {
			matches = append(matches, CarveoutLineMatch{BeforeIndex: beforeIndex, AfterIndex: afterIndex})
		}
	}
	if !matcher.check() {
		return matcher.failure(prepared, options.Fallback)
	}
	if options.RequireMatch && len(matches) == 0 {
		if !matcher.check() {
			return matcher.failure(prepared, options.Fallback)
		}
		return matchingFailure(prepared, CarveoutMatchImpossibility, CarveoutMatchRequiredMatchUnavailable, matcher.inspections, options.Fallback)
	}
	if !matcher.check() {
		return matcher.failure(prepared, options.Fallback)
	}
	return CarveoutMatchAnalysis{State: CarveoutMatchCompleted, Result: &CarveoutMatchResult{moved: len(matches), matches: matches}, Provenance: CarveoutMatchProvenance{State: CarveoutMatchCompleted, PairInspections: matcher.inspections}}
}
func (m *preparedCarveoutMatcher) scan(prepared PreparedCarveout, consumed []bool, selected []int, fuzzy bool) bool {
	for afterIndex := range selected {
		if selected[afterIndex] >= 0 {
			continue
		}
		if !m.check() {
			return false
		}
		afterLine, _ := prepared.AfterLine(afterIndex)
		for beforeIndex := 0; beforeIndex < len(consumed); beforeIndex++ {
			if consumed[beforeIndex] {
				continue
			}
			if !m.inspect() {
				return false
			}
			beforeLine, _ := prepared.BeforeLine(beforeIndex)
			matched := false
			if fuzzy {
				var completed bool
				matched, completed = m.similar(beforeLine, afterLine)
				if !completed {
					return false
				}
			} else {
				matched = beforeLine == afterLine
				if !m.check() {
					return false
				}
			}
			if matched {
				consumed[beforeIndex], selected[afterIndex] = true, beforeIndex
				break
			}
		}
	}
	return true
}
func matchingPrecondition(options CarveoutMatcherOptions) CarveoutMatchReason {
	switch {
	case options.PairCap <= 0:
		return CarveoutMatchInvalidPairCap
	case options.TimeCap <= 0:
		return CarveoutMatchInvalidTimeCap
	case options.Clock == nil:
		return CarveoutMatchMissingClock
	}
	return ""
}
func (m *preparedCarveoutMatcher) check() bool {
	if !m.clock.Now().Before(m.deadline) {
		m.reason = CarveoutMatchDeadlineExceeded
		return false
	}
	return true
}
func (m *preparedCarveoutMatcher) inspect() bool {
	if !m.check() {
		return false
	}
	if m.inspections == m.pairCap {
		m.reason = CarveoutMatchPairCapExceeded
		return false
	}
	m.inspections++
	return true
}
func (m *preparedCarveoutMatcher) similar(before, after string) (bool, bool) {
	if !m.check() {
		return false, false
	}
	maximum := len(before)
	if len(after) > maximum {
		maximum = len(after)
	}
	if !m.check() {
		return false, false
	}
	previous, current := make([]int, len(after)+1), make([]int, len(after)+1)
	if !m.check() {
		return false, false
	}
	for index := range previous {
		if !m.check() {
			return false, false
		}
		previous[index] = index
	}
	for beforeIndex := 1; beforeIndex <= len(before); beforeIndex++ {
		if !m.check() {
			return false, false
		}
		current[0] = beforeIndex
		for afterIndex := 1; afterIndex <= len(after); afterIndex++ {
			if !m.check() {
				return false, false
			}
			cost := 0
			if before[beforeIndex-1] != after[afterIndex-1] {
				cost = 1
			}
			current[afterIndex] = min(previous[afterIndex]+1, current[afterIndex-1]+1, previous[afterIndex-1]+cost)
		}
		previous, current = current, previous
	}
	if !m.check() {
		return false, false
	}
	return previous[len(after)]*100 <= maximum*20, true
}
func (m *preparedCarveoutMatcher) failure(prepared PreparedCarveout, fallback *CarveoutMatchFallback) CarveoutMatchAnalysis {
	return matchingFailure(prepared, CarveoutMatchTimeout, m.reason, m.inspections, fallback)
}
func matchingFailure(prepared PreparedCarveout, state CarveoutMatchState, reason CarveoutMatchReason, inspections int, fallback *CarveoutMatchFallback) CarveoutMatchAnalysis {
	analysis := CarveoutMatchAnalysis{State: state, Provenance: CarveoutMatchProvenance{State: state, Reason: reason, PairInspections: inspections}}
	if fallback != nil && fallback.Audited && strings.TrimSpace(fallback.Reason) != "" && strings.TrimSpace(fallback.Evidence) != "" && fallback.Moved >= 0 && fallback.Moved <= prepared.BeforeLineCount() && fallback.Moved <= prepared.AfterLineCount() {
		retained := *fallback
		analysis.State, analysis.Fallback = CarveoutMatchFallbackApplied, &retained
	}
	return analysis
}
