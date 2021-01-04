package hcledit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Writer interface {
	WriteFile(path string) error

	Write(w io.Writer) error

	OverWriteFile() error

	// Bytes returns a buffer containing the source code resulting from the
	// tokens underlying the receiving file. If any updates have been made via
	// the AST API, these will be reflected in the result.
	Bytes() []byte
}

func (h *hclEditImpl) WriteFile(path string) error {
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

func (h *hclEditImpl) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s", h.Bytes())
	return err
}

func (h *hclEditImpl) OverWriteFile() error {
	if h.path == "" {
		return fmt.Errorf("OverWriteFile can be used only when you create editor via ReadFile()")
	}

	return h.WriteFile(h.path)
}

func (h *hclEditImpl) Bytes() []byte {
	return h.writeFile.Bytes()
}
