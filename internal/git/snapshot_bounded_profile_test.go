package gitadapter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/internal/w1w2/contract"
)

func fixedProfileOutputs(count, rawSize int) [][]byte {
	prefix := ":000000 100644 " + strings.Repeat("0", 40) + " " + strings.Repeat("c", 40) + " A\x00"
	var raw, numstat bytes.Buffer
	for i := 0; i < count; i++ {
		path := fmt.Sprintf("p%05d", i)
		if rawSize != 0 {
			payload := rawSize - count*(len(prefix)+1)
			length := payload / count
			if i < payload%count {
				length++
			}
			for len(path) < length {
				left := length - len(path)
				if left == 1 {
					path += "x"
				} else {
					path += "/" + strings.Repeat("x", min(255, left-1))
				}
			}
		}
		raw.WriteString(prefix + path + "\x00")
		numstat.WriteString("1\t0\t" + path + "\x00")
	}
	return [][]byte{[]byte("sha1\n"), []byte(strings.Repeat("a", 40) + "\n"), nil, []byte(strings.Repeat("b", 40) + "\n"), nil, raw.Bytes(), numstat.Bytes(), nil, nil}
}

func padFixedProfileTotal(outputs [][]byte, total int) {
	remaining := total
	for _, out := range outputs {
		remaining -= len(out)
	}
	// Whitespace is valid around parsed metadata; successful cat-file -e output
	// is ignored but still charged. No individual response crosses its cap.
	for _, index := range []int{0, 1, 2, 3} {
		padding := min(remaining, contract.MaxGitResponseBytes-len(outputs[index]))
		outputs[index] = append(outputs[index], bytes.Repeat([]byte(" "), padding)...)
		remaining -= padding
	}
	if remaining != 0 {
		panic("invalid bounded aggregate fixture")
	}
}

func TestBoundedSnapshotFixedProfileBoundaries(t *testing.T) {
	bounds := SnapshotBounds{contract.MaxGitResponseBytes, contract.MaxGitBytes, contract.MaxDiffEntries}
	for _, surface := range []string{"response", "aggregate", "entries"} {
		for _, delta := range []int{-1, 0, 1} {
			t.Run(fmt.Sprintf("%s/%+d", surface, delta), func(t *testing.T) {
				count, size := 1, 0
				if surface == "response" {
					count, size = bounds.MaxEntries, bounds.MaxResponseBytes+delta
				} else if surface == "entries" {
					count = bounds.MaxEntries + delta
				}
				outputs := fixedProfileOutputs(count, size)
				if surface == "aggregate" {
					padFixedProfileTotal(outputs, bounds.MaxTotalBytes+delta)
				}
				actual := len(outputs[5])
				if surface == "response" && actual != bounds.MaxResponseBytes+delta {
					t.Fatalf("response fixture size %d", actual)
				}
				if surface == "aggregate" {
					actual = 0
					for _, out := range outputs {
						if len(out) > bounds.MaxResponseBytes {
							t.Fatal("aggregate fixture exceeds individual cap")
						}
						actual += len(out)
					}
					if actual != bounds.MaxTotalBytes+delta {
						t.Fatalf("aggregate fixture size %d", actual)
					}
				}
				calls, remaining := 0, bounds.MaxTotalBytes
				runner := func(_ context.Context, _ string, intCap int, _ ...string) ([]byte, error) {
					if intCap != min(bounds.MaxResponseBytes, remaining) || calls >= len(outputs) {
						t.Fatal("unexpected acquisition budget or command count")
					}
					out := outputs[calls]
					calls++
					remaining -= len(out)
					return out, nil
				}
				got, err := AcquireBoundedWithRunner(context.Background(), runner, Request{"/fixture", "base", "head"}, bounds)
				if delta > 0 {
					if !errors.Is(err, ErrSnapshotResourceLimit) || !reflect.DeepEqual(got, evidence.CommittedSnapshot{}) {
						t.Fatalf("overrun returned partial/success: entries=%d error=%v", len(got.Entries()), err)
					}
				} else if err != nil || len(got.Entries()) != count || calls != len(outputs) {
					t.Fatalf("admitted fixture: entries=%d calls=%d error=%v", len(got.Entries()), calls, err)
				}
			})
		}
	}
}
