package command

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRunUpdate(t *testing.T) {
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
		"e2-highmem-2",
		filename,
	}

	if err := runUpdate(&UpdateOptions{Type: "string"}, args); err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	want := `
resource "google_container_node_pool" "nodes1" {
  node_config {
    preemptible  = false
    machine_type = "e2-highmem-2"
  }
}
`
	diff := cmp.Diff(want, string(got), cmpopts.AcyclicTransformer("multiline", func(s string) []string {
		return strings.Split(s, "\n")
	}))

	if diff != "" {
		t.Fatalf("Update mismatch (-want +got):\n%s", diff)
	}

}
