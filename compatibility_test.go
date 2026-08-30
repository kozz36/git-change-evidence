package changeevidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestImportedPredecessorCorpus(t *testing.T) {
	const manifest = `{"authored_by":"git-change-evidence successor","predecessor":{"commit":"a0fc7b26ff8a0e0a61baa586b32c46841611c806","files":[{"blob":"ff9c62614cdf1e7891eccad01386bce13500b9e0","local":"contracts-v1.json","path":"tests/fixtures/sd7/compatibility/v1/contracts-v1.json","sha256":"84bba5b7b1471a1c12477de6a46a2f6f175fd09bdd175ecaf81c698902776ca5"},{"blob":"1dc4109145e697523e82c4199fd64062babf9d0b","local":"accounting-v1.json","path":"tests/fixtures/sd7/compatibility/v1/accounting-v1.json","sha256":"f265bc48750ccd68e41dcb64ec68134b9fb7db69655e2c215e39f8c6c8794c3f"},{"blob":"64fb760df29cb1c7f0fa74866779213b2c1616f2","local":"manifest.sha256","path":"tests/fixtures/sd7/compatibility/v1/manifest.sha256","sha256":"f8db149ef966e1eeab42e36e42e299ba4b197baeac22bc75cfb4b9cd55b8caac"}],"repository":"CNSIC"},"schema":"git-change-evidence.compatibility-manifest/v1"}`
	root := filepath.Join("testdata", "compatibility", "v1")
	got, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte(manifest)) {
		t.Fatalf("manifest bytes = %q, want %q", got, manifest)
	}
	for _, file := range []struct{ local, sha256 string }{
		{"contracts-v1.json", "84bba5b7b1471a1c12477de6a46a2f6f175fd09bdd175ecaf81c698902776ca5"},
		{"accounting-v1.json", "f265bc48750ccd68e41dcb64ec68134b9fb7db69655e2c215e39f8c6c8794c3f"},
		{"manifest.sha256", "f8db149ef966e1eeab42e36e42e299ba4b197baeac22bc75cfb4b9cd55b8caac"},
	} {
		contents, err := os.ReadFile(filepath.Join(root, file.local))
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(contents)
		if got := hex.EncodeToString(sum[:]); got != file.sha256 {
			t.Fatalf("%s SHA-256 = %s, want %s", file.local, got, file.sha256)
		}
	}
}
