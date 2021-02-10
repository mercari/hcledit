package hcledit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WriteFile writes the new contents to the given file, creating it if it does
// not exist.
func (h *HCLEditor) WriteFile(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return h.Write(f)
}

// Write writes the new contents to the given io.Writer.
func (h *HCLEditor) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s", h.Bytes())
	return err
}

// OverWriteFile writes the new contents to the file that has first been read
// via ReadFile.
func (h *HCLEditor) OverWriteFile() error {
	if h.path == "" {
		return fmt.Errorf("OverWriteFile can be used only when you create editor via ReadFile()")
	}

	return h.WriteFile(h.path)
}

// Bytes returns a buffer containing the source code resulting from the
// tokens underlying the receiving file. If any updates have been made via
// the AST API, these will be reflected in the result.
func (h *HCLEditor) Bytes() []byte {
	return h.writeFile.Bytes()
}
