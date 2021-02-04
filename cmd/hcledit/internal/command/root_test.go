package command

import (
	"io/ioutil"
	"testing"
)

func tempFile(t *testing.T, contents string) string {
	t.Helper()
	f, err := ioutil.TempFile(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := f.WriteString(contents); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func readFile(t *testing.T, filename string) string {
	t.Helper()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	return string(b)
}
