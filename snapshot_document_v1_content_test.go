package changeevidence_test

import (
	"bytes"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
)

func TestMeasureSnapshotContentV1(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  evidence.SnapshotContentV1
	}{
		{"nil", nil, evidence.SnapshotContentV1{ByteLength: 0, ContainsNUL: false, LogicalLines: 0}},
		{"empty", []byte{}, evidence.SnapshotContentV1{ByteLength: 0, ContainsNUL: false, LogicalLines: 0}},
		{"single_lf", []byte("\n"), evidence.SnapshotContentV1{ByteLength: 1, ContainsNUL: false, LogicalLines: 1}},
		{"two_lf", []byte("\n\n"), evidence.SnapshotContentV1{ByteLength: 2, ContainsNUL: false, LogicalLines: 2}},
		{"terminated", []byte("a\n"), evidence.SnapshotContentV1{ByteLength: 2, ContainsNUL: false, LogicalLines: 1}},
		{"unterminated", []byte("a"), evidence.SnapshotContentV1{ByteLength: 1, ContainsNUL: false, LogicalLines: 1}},
		{"unterminated_tail", []byte("a\nb"), evidence.SnapshotContentV1{ByteLength: 3, ContainsNUL: false, LogicalLines: 2}},
		{"crlf", []byte("\r\n"), evidence.SnapshotContentV1{ByteLength: 2, ContainsNUL: false, LogicalLines: 1}},
		{"crlf_tail", []byte("a\r\nb"), evidence.SnapshotContentV1{ByteLength: 4, ContainsNUL: false, LogicalLines: 2}},
		{"cr_is_not_separator", []byte("a\rb"), evidence.SnapshotContentV1{ByteLength: 3, ContainsNUL: false, LogicalLines: 1}},
		{"nul", []byte{0}, evidence.SnapshotContentV1{ByteLength: 1, ContainsNUL: true, LogicalLines: 1}},
		{"nul_before_lf", []byte("a\x00\n"), evidence.SnapshotContentV1{ByteLength: 3, ContainsNUL: true, LogicalLines: 1}},
		{"nul_in_tail", []byte("a\n\x00"), evidence.SnapshotContentV1{ByteLength: 3, ContainsNUL: true, LogicalLines: 2}},
		{"non_utf8_bytes", []byte{0xff, '\n', 0xfe}, evidence.SnapshotContentV1{ByteLength: 3, ContainsNUL: false, LogicalLines: 2}},
		{"utf8_is_measured_as_bytes", []byte("é\n"), evidence.SnapshotContentV1{ByteLength: 3, ContainsNUL: false, LogicalLines: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := append([]byte(nil), tt.input...)
			got := evidence.MeasureSnapshotContentV1(tt.input)
			if got != tt.want {
				t.Fatalf("metrics = %#v, want %#v", got, tt.want)
			}
			if !bytes.Equal(tt.input, before) {
				t.Fatalf("input mutated: got %x, want %x", tt.input, before)
			}
		})
	}
}

func TestMeasureSnapshotContentV1ValueIsolation(t *testing.T) {
	input := []byte("a\x00\nb")
	first := evidence.MeasureSnapshotContentV1(input)
	first.ByteLength = 99
	first.ContainsNUL = false
	first.LogicalLines = 99
	want := evidence.SnapshotContentV1{ByteLength: 4, ContainsNUL: true, LogicalLines: 2}
	if got := evidence.MeasureSnapshotContentV1(input); got != want {
		t.Fatalf("fresh metrics = %#v, want %#v; prior value = %#v", got, want, first)
	}
}
