package hcledit

import (
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// New constructs a new HCL file with no content which is ready to be mutated.
func New() (*HCLEditor, error) {
	return &HCLEditor{
		writeFile: hclwrite.NewEmptyFile(),
	}, nil
}

// ReadFile reads HCL file in the given path and returns operation interface for it.
func ReadFile(path string) (*HCLEditor, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	editor, err := Read(f, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	editor.path = path

	return editor, err
}

// Read reads HCL file from the given io.Reader and returns operation interface for it.
func Read(r io.Reader, filename string) (*HCLEditor, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	writeFile, diags := hclwrite.ParseConfig(buf, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	return &HCLEditor{
		filename:  filename,
		writeFile: writeFile,
	}, nil
}
