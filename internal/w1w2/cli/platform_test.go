package cli

import (
	"bytes"
	"context"
	"runtime"
	"testing"
)

func TestUnsupportedPlatformBeforeInputAcquisition(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("requires an actual non-Linux runtime")
	}
	for _, args := range [][]string{
		{"inventory", "--source-root", "/fixture", "--policy", "-", "--paths", "missing-selection"},
		{"result", "--repo", "/fixture", "--base", "base", "--head", "head", "--policy", "-", "--inventory", "missing-inventory"},
		{"report", "--policy", "-", "--inventory", "missing-inventory", "--result", "missing-result"},
	} {
		t.Run(args[0], func(t *testing.T) {
			var out, diagnostic bytes.Buffer
			code := Run(context.Background(), args, unusedInput{}, &out, &diagnostic)
			if code != 3 || out.Len() != 0 || diagnostic.String() != "w1w2 "+args[0]+": unsupported_platform\n" {
				t.Fatalf("code=%d stdout=%q stderr=%q", code, out.Bytes(), diagnostic.String())
			}
		})
	}
}
