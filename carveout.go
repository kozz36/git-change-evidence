package changeevidence

import (
	"crypto/sha256"
	"encoding/hex"
)

const (
	CarveoutSafeMaxBytes     = 16 << 20
	CarveoutSafeMaxLines     = 1 << 18
	CarveoutSafeMaxLineBytes = 1 << 20
)

type CarveoutLimit string

const (
	CarveoutByteLimit     CarveoutLimit = "byte_cap"
	CarveoutLineLimit     CarveoutLimit = "line_cap"
	CarveoutLineByteLimit CarveoutLimit = "max_line_bytes"
)

type CarveoutState string

const (
	CarveoutInvalidBound               CarveoutState = "invalid_bound"
	CarveoutUnavailableRequiredContent CarveoutState = "unavailable_required_content"
	CarveoutResourceBoundExceeded      CarveoutState = "resource_bound_exceeded"
	CarveoutPrepared                   CarveoutState = "prepared"
)

type CarveoutContentName string

const (
	CarveoutBeforeContent CarveoutContentName = "before"
	CarveoutAfterContent  CarveoutContentName = "after"
)

type CarveoutBounds struct{ ByteCap, LineCap, MaxLineBytes int }
type CarveoutContent struct {
	Available bool
	Raw       string
}
type CarveoutInput struct{ Before, After CarveoutContent }

type PreparedCarveout struct {
	before, after                   []string
	beforeRawDigest, afterRawDigest Digest
}

func (p PreparedCarveout) BeforeLineCount() int { return len(p.before) }
func (p PreparedCarveout) AfterLineCount() int  { return len(p.after) }
func (p PreparedCarveout) BeforeLine(index int) (string, bool) {
	return carveoutLineAt(p.before, index)
}
func (p PreparedCarveout) AfterLine(index int) (string, bool) { return carveoutLineAt(p.after, index) }
func (p PreparedCarveout) BeforeRawDigest() Digest            { return p.beforeRawDigest }
func (p PreparedCarveout) AfterRawDigest() Digest             { return p.afterRawDigest }

func carveoutLineAt(lines []string, index int) (string, bool) {
	if index < 0 || index >= len(lines) {
		return "", false
	}
	return lines[index], true
}

type CarveoutPreparation struct {
	State    CarveoutState
	Limit    CarveoutLimit
	Content  CarveoutContentName
	Prepared PreparedCarveout
}

func PrepareCarveout(input CarveoutInput, bounds CarveoutBounds) CarveoutPreparation {
	if limit := invalidCarveoutBound(bounds); limit != "" {
		return CarveoutPreparation{State: CarveoutInvalidBound, Limit: limit}
	}
	if !input.Before.Available {
		return CarveoutPreparation{State: CarveoutUnavailableRequiredContent, Content: CarveoutBeforeContent}
	}
	if !input.After.Available {
		return CarveoutPreparation{State: CarveoutUnavailableRequiredContent, Content: CarveoutAfterContent}
	}
	if !withinCarveoutByteCap(len(input.Before.Raw), len(input.After.Raw), bounds.ByteCap) {
		return CarveoutPreparation{State: CarveoutResourceBoundExceeded, Limit: CarveoutByteLimit}
	}
	beforeCount, limit := censusCarveoutLines(input.Before.Raw, 0, bounds)
	if limit != "" {
		return CarveoutPreparation{State: CarveoutResourceBoundExceeded, Limit: limit}
	}
	afterCount, limit := censusCarveoutLines(input.After.Raw, beforeCount, bounds)
	if limit != "" {
		return CarveoutPreparation{State: CarveoutResourceBoundExceeded, Limit: limit}
	}
	return CarveoutPreparation{State: CarveoutPrepared, Prepared: PreparedCarveout{
		before: ownCarveoutLines(input.Before.Raw, beforeCount), after: ownCarveoutLines(input.After.Raw, afterCount-beforeCount),
		beforeRawDigest: digestCarveoutRaw(input.Before.Raw), afterRawDigest: digestCarveoutRaw(input.After.Raw),
	}}
}

func invalidCarveoutBound(bounds CarveoutBounds) CarveoutLimit {
	switch {
	case bounds.ByteCap <= 0 || bounds.ByteCap > CarveoutSafeMaxBytes:
		return CarveoutByteLimit
	case bounds.LineCap <= 0 || bounds.LineCap > CarveoutSafeMaxLines:
		return CarveoutLineLimit
	case bounds.MaxLineBytes <= 0 || bounds.MaxLineBytes > CarveoutSafeMaxLineBytes:
		return CarveoutLineByteLimit
	}
	return ""
}

func withinCarveoutByteCap(left, right, limit int) bool { return left <= limit && right <= limit-left }

func censusCarveoutLines(raw string, count int, bounds CarveoutBounds) (int, CarveoutLimit) {
	start := 0
	for index := 0; index < len(raw); index++ {
		if raw[index] != '\n' {
			continue
		}
		length := index - start
		if length > 0 && raw[index-1] == '\r' {
			length--
		}
		if length > 0 {
			if length > bounds.MaxLineBytes {
				return 0, CarveoutLineByteLimit
			}
			if count == bounds.LineCap {
				return 0, CarveoutLineLimit
			}
			count++
		}
		start = index + 1
	}
	if start < len(raw) {
		if len(raw)-start > bounds.MaxLineBytes {
			return 0, CarveoutLineByteLimit
		}
		if count == bounds.LineCap {
			return 0, CarveoutLineLimit
		}
		count++
	}
	return count, ""
}

func ownCarveoutLines(raw string, count int) []string {
	lines, start, line := make([]string, count), 0, 0
	for index := 0; index < len(raw); index++ {
		if raw[index] != '\n' {
			continue
		}
		end := index
		if end > start && raw[end-1] == '\r' {
			end--
		}
		if end > start {
			lines[line], line = raw[start:end], line+1
		}
		start = index + 1
	}
	if start < len(raw) {
		lines[line] = raw[start:]
	}
	return lines
}

func digestCarveoutRaw(raw string) Digest {
	sum := sha256.Sum256([]byte(raw))
	return Digest(hex.EncodeToString(sum[:]))
}
