package gitadapter_test

import (
	"bytes"
	"errors"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	gitadapter "github.com/kozz36/git-change-evidence/internal/git"
)

func TestParseSnapshotMetadataRejectsInvalidFramesAndPairing(t *testing.T) {
	a, b := metadataRawFixture("M", "a", "", 40), metadataRawFixture("M", "b", "", 40)
	x, y := metadataNumstatFixture("2", "1", "a", ""), metadataNumstatFixture("2", "1", "b", "")
	pairRaw := append(append([]byte(nil), a...), b...)
	pairNum := append(append([]byte(nil), x...), y...)
	rename := metadataRawFixture("R100", "new", "old", 40)
	cases := []struct {
		name         string
		raw, numstat []byte
	}{
		{"missing_numstat", a, nil}, {"extra_numstat", a, pairNum},
		{"duplicate_numstat_hidden_by_map", pairRaw, append(append([]byte(nil), pairNum...), x...)},
		{"duplicate_raw_with_unused_extra", append(append([]byte(nil), a...), a...), pairNum},
		{"duplicate_pair", append(append([]byte(nil), rename...), rename...), append(metadataNumstatFixture("1", "0", "new", "old"), metadataNumstatFixture("1", "0", "new", "old")...)},
		{"wrong_pair", rename, metadataNumstatFixture("1", "0", "new", "other")},
		{"pair_not_single_path", rename, metadataNumstatFixture("1", "0", "new", "")},
		{"no_path_normalization", a, metadataNumstatFixture("2", "1", "./a", "")},
		{"raw_unterminated", a[:len(a)-1], x}, {"numstat_unterminated", a, x[:len(x)-1]},
		{"raw_empty_record", append(append([]byte(nil), a...), 0), x},
		{"numstat_empty_record", a, append(append([]byte(nil), x...), 0)},
		{"missing_rename_path", bytes.TrimSuffix(rename, []byte("new\x00")), metadataNumstatFixture("1", "0", "new", "old")},
		{"unknown_status", metadataRawFixture("U", "a", "", 40), x},
		{"score_over_100", metadataRawFixture("R101", "new", "old", 40), metadataNumstatFixture("1", "0", "new", "old")},
		{"score_too_long", metadataRawFixture("R0000", "new", "old", 40), metadataNumstatFixture("1", "0", "new", "old")},
		{"score_not_decimal", metadataRawFixture("Rfoo", "new", "old", 40), metadataNumstatFixture("1", "0", "new", "old")},
		{"score_on_addition", metadataRawFixture("A100", "a", "", 40), x},
		{"only_numstat", nil, x},
		{"missing_numstat_rename_path", rename, []byte("1\t0\t\x00old\x00")},
		{"malformed_header", []byte(":broken\x00a\x00"), x},
		{"numstat_fields", a, []byte("2\ta\x00")},
		{"mixed_binary_added", a, metadataNumstatFixture("-", "1", "a", "")},
		{"mixed_binary_deleted", a, metadataNumstatFixture("2", "-", "a", "")},
	}
	for _, token := range []string{"", "01", "+1", "-1", "1e0", "1.0", "true", " 1", "١", "18446744073709551616"} {
		cases = append(cases, struct {
			name         string
			raw, numstat []byte
		}{"added_token_" + token, a, metadataNumstatFixture(token, "0", "a", "")})
		cases = append(cases, struct {
			name         string
			raw, numstat []byte
		}{"deleted_token_" + token, a, metadataNumstatFixture("0", token, "a", "")})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw, numstat := append([]byte(nil), tc.raw...), append([]byte(nil), tc.numstat...)
			rawBefore, numstatBefore := append([]byte(nil), raw...), append([]byte(nil), numstat...)
			got, err := gitadapter.ParseSnapshotMetadata(raw, numstat)
			var snapshotError *evidence.SnapshotError
			if !errors.As(err, &snapshotError) || snapshotError.Code != evidence.SnapshotRacing || err.Error() != "snapshot acquisition: racing_object" {
				t.Fatalf("unexpected error = %v", err)
			}
			if got != nil || !bytes.Equal(raw, rawBefore) || !bytes.Equal(numstat, numstatBefore) {
				t.Fatal("rejection leaked partial records or mutated input")
			}
		})
	}
}
