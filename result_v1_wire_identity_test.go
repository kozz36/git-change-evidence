package changeevidence

import (
	"bytes"
	"testing"
)

func TestResultV1CanonicalWireUsesClosedVariants(t *testing.T) {
	document, err := wireResultDocument(
		[]ResultEntryV1{
			{Path: []byte("p"), Category: "source", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1}},
			{Path: []byte("a"), Category: "other", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryNonCountableV1, Additions: 9, Deletions: 8}},
		},
		[]ResultCategoryTotalV1{{Category: "source"}, {Category: "other", NonCountable: 1}},
		[]ResultObservationV1{
			{Kind: ResultLineThresholdObservationV1, Category: "source", LineThresholdLimit: 7, Available: true},
			{Kind: ResultLineThresholdObservationV1, Category: "other", LineThresholdLimit: 3, Actual: 9, Exceeded: true},
			{Kind: ResultRatioObservationV1, Category: "source", Reference: "other", RatioLimit: "0.5", Available: true, Denominator: 1},
			{Kind: ResultRatioObservationV1, Category: "other", Reference: "source", RatioLimit: "1", Numerator: 9, Denominator: 8, Exceeded: true},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`{"schema":"git-change-evidence.accounting-result/v1","accounting_policy_sha256":"` + string(wirePolicyDigest) + `","inventory_sha256":"` + string(wireInventoryDigest) + `","revisions":{"base":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","head":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},"entries":[{"path_b64":"cA==","category":"source","measurement":{"kind":"countable","additions":0,"deletions":0}},{"path_b64":"YQ==","category":"other","measurement":{"kind":"non_countable"}}],"totals":[{"category":"source","additions":0,"deletions":0,"non_countable":0},{"category":"other","additions":0,"deletions":0,"non_countable":1}],"observations":[{"kind":"line_threshold","category":"source","limit":7,"available":true,"actual":0,"exceeded":false},{"kind":"line_threshold","category":"other","limit":3,"available":false,"exceeded":false},{"kind":"ratio","category":"source","reference":"other","limit":"0.5","available":true,"numerator":0,"denominator":1,"exceeded":false},{"kind":"ratio","category":"other","reference":"source","limit":"1","available":false,"exceeded":false}]}` + "\n")
	if got := document.CanonicalBytes(); !bytes.Equal(got, want) {
		t.Fatalf("canonical bytes = %q, want %q", got, want)
	}
}

func TestResultV1CanonicalWirePreservesPathByteIdentity(t *testing.T) {
	for _, test := range []struct {
		name, leftB64, rightB64 string
		left, right             []byte
	}{
		{"ff and fe", "/w==", "/g==", []byte{0xff}, []byte{0xfe}},
		{"LF and CRLF adjacent", "YQpi", "YQ0KYg==", []byte("a\nb"), []byte("a\r\nb")},
		{"padding required", "YQ==", "YWI=", []byte("a"), []byte("ab")},
	} {
		t.Run(test.name, func(t *testing.T) {
			left, err := wireResultDocument([]ResultEntryV1{{Path: test.left, Category: "source", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1}}}, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			right, err := wireResultDocument([]ResultEntryV1{{Path: test.right, Category: "source", Measurement: ResultEntryMeasurementV1{Kind: ResultEntryCountableV1}}}, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Equal(left.CanonicalBytes(), right.CanonicalBytes()) || left.Digest() == right.Digest() {
				t.Fatal("distinct raw paths share canonical identity")
			}
			for _, encoded := range []string{test.leftB64, test.rightB64} {
				if !bytes.Contains(append(left.CanonicalBytes(), right.CanonicalBytes()...), []byte(`"path_b64":"`+encoded+`"`)) {
					t.Fatalf("canonical bytes omit padded Base64 %q", encoded)
				}
			}
		})
	}
}
