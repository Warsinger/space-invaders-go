//go:build !js

package score

import (
	"os"
	"testing"
)

func TestNativeRoundTrip(t *testing.T) {
	t.Chdir(t.TempDir())
	if Load() != 0 {
		t.Fatal("missing score must start at zero")
	}
	if err := Save(125); err != nil {
		t.Fatal(err)
	}
	if Load() != 125 {
		t.Fatal("score did not survive reload")
	}
	if err := os.WriteFile(File, []byte("corrupt"), 0644); err != nil {
		t.Fatal(err)
	}
	if Load() != 0 {
		t.Fatal("corrupt score must start at zero")
	}
}

func TestSaveFailure(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile("score", []byte("not a directory"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Save(1); err == nil {
		t.Fatal("expected failure when score directory is a file")
	}
}
