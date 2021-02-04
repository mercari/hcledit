package hcledit_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mercari/hcledit"
)

func TestWriteFile(t *testing.T) {
	originalContent := `
attribute1 = "str1"
`
	editor, err := hcledit.Read(strings.NewReader(originalContent), "")
	if err != nil {
		t.Fatal(err)
	}

	if err := editor.Create("attribute2", "C", hcledit.WithAfter("attribute1")); err != nil {
		t.Fatal(err)
	}

	if err := editor.Update("attribute1", "U"); err != nil {
		t.Fatal(err)
	}

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.hcl")
	if err := editor.WriteFile(tempFile); err != nil {
		t.Fatal(err)
	}

	want := []byte(`
attribute1 = "U"
attribute2 = "C"
`)
	got, err := ioutil.ReadFile(tempFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("WriteFile() mismatch:\ngot:%s\nwant:%s\n", got, want)
	}
}

func TestOverWriteFile(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.hcl")
	if err := ioutil.WriteFile(tempFile, []byte(`
attribute1 = "str1"
`), 0600); err != nil {
		t.Fatal(err)
	}

	editor, err := hcledit.ReadFile(tempFile)
	if err != nil {
		t.Fatal(err)
	}

	if err := editor.Create("attribute2", "C", hcledit.WithAfter("attribute1")); err != nil {
		t.Fatal(err)
	}

	if err := editor.Update("attribute1", "U"); err != nil {
		t.Fatal(err)
	}

	if err := editor.OverWriteFile(); err != nil {
		t.Fatal(err)
	}

	want := []byte(`
attribute1 = "U"
attribute2 = "C"
`)
	got, err := ioutil.ReadFile(tempFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("WriteFile() mismatch:\ngot:%s\nwant:%s\n", got, want)
	}
}
