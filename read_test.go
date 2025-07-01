package hcledit_test

import (
	"testing"

	"go.mercari.io/hcledit"
)

func TestReadFile(t *testing.T) {
	editor, err := hcledit.ReadFile("cmd/hcledit/internal/command/fixture/file.tf")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if editor == nil {
		t.Fatalf("editor should not be nil")
	}
}
