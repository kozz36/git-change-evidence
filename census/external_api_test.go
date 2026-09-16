package census_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	evidence "github.com/kozz36/git-change-evidence"
	"github.com/kozz36/git-change-evidence/census"
)

func TestExternalValueAPIAndOwnership(t *testing.T) {
	path := []byte("external.go")
	content := []byte("package p\nfunc f(){ (os.Open)(\"x\") }\n")
	policy, err := evidence.NewAccountingPolicyV1(evidence.PolicyInput{
		Categories: []evidence.CategoryInput{{Name: "go"}},
		Default:    "go",
	})
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := evidence.NewInventoryV1(policy, []evidence.InventoryEntryInput{{Path: path, Content: content}})
	if err != nil {
		t.Fatal(err)
	}
	result, err := census.GoASTV1(inventory, []census.File{{Path: path, Content: content}}, census.QueryV1{
		Version: census.SelectorCallQueryV1, Receiver: "os", Selector: "Open",
	}, census.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	if result.InventoryDigest() != inventory.Digest() || result.Query() != (census.QueryV1{Version: census.SelectorCallQueryV1, Receiver: "os", Selector: "Open"}) || result.ExtractorVersion() != census.ExtractorVersion {
		t.Fatalf("bindings = (%q, %#v, %q)", result.InventoryDigest(), result.Query(), result.ExtractorVersion())
	}
	matches := result.Matches()
	if len(matches) != 1 || !bytes.Equal(matches[0].Path, []byte("external.go")) || matches[0].Call.Start != (census.Position{Offset: 20, Line: 2, Column: 11}) || matches[0].Call.End != (census.Position{Offset: 34, Line: 2, Column: 25}) || matches[0].Selector.Start != (census.Position{Offset: 21, Line: 2, Column: 12}) || matches[0].Selector.End != (census.Position{Offset: 28, Line: 2, Column: 19}) {
		t.Fatalf("match = %#v", matches)
	}
	fileSum := sha256.Sum256(content)
	fragmentSum := sha256.Sum256([]byte("(os.Open)(\"x\")"))
	if matches[0].FileSHA256 != evidence.Digest(hex.EncodeToString(fileSum[:])) || matches[0].FragmentSHA256 != evidence.Digest(hex.EncodeToString(fragmentSum[:])) {
		t.Fatalf("hashes = %#v", matches[0])
	}
	path[0], content[0], matches[0].Path[0] = 'X', 'X', 'X'
	matches = append(matches, census.MatchV1{})
	again := result.Matches()
	if len(again) != 1 || !bytes.Equal(again[0].Path, []byte("external.go")) || again[0].Call.Start.Offset != 20 {
		t.Fatalf("mutable data escaped: %#v", again)
	}
}
