package command

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRunDelete(t *testing.T) {
	filename := tempFile(t, `
resource "google_container_node_pool" "nodes1" {
  node_config {
    preemptible  = false
    machine_type = "e2-medium"
  }
}
`)

	args := []string{
		"resource.google_container_node_pool.*.node_config.machine_type",
		filename,
	}
	if err := runDelete(args); err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	want := `
resource "google_container_node_pool" "nodes1" {
  node_config {
    preemptible = false
  }
}
`
	diff := cmp.Diff(want, string(got), cmpopts.AcyclicTransformer("multiline", func(s string) []string {
		return strings.Split(s, "\n")
	}))

	if diff != "" {
		t.Fatalf("Delete mismatch (-want +got):\n%s", diff)
	}
}
