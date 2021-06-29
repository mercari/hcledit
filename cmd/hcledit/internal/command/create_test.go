package command

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRunCreate(t *testing.T) {
	cases := map[string]struct {
		opts *CreateOptions
		want string
	}{
		"WithoutAdditionalOptions": {
			opts: &CreateOptions{},
			want: `resource "google_container_node_pool" "nodes1" {
  node_config {
    preemptible  = false
    machine_type = "e2-medium"
    disk_size_gb = "100"
  }
}
`,
		},
		"WithOptionWithAfter": {
			opts: &CreateOptions{
				Type:  "string",
				After: "preemptible",
			},
			want: `resource "google_container_node_pool" "nodes1" {
  node_config {
    preemptible  = false
    disk_size_gb = "100"
    machine_type = "e2-medium"
  }
}
`,
		},
		"WithOptionComment": {
			opts: &CreateOptions{
				Comment: "// TODO: Testing",
			},
			want: `resource "google_container_node_pool" "nodes1" {
  node_config {
    preemptible  = false
    machine_type = "e2-medium"
    // TODO: Testing
    disk_size_gb = "100"
  }
}
`,
		},
	}

	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			filename := tempFile(t, `resource "google_container_node_pool" "nodes1" {
  node_config {
    preemptible  = false
    machine_type = "e2-medium"
  }
}
`)

			tc.opts.Type = "string" // ensure default value

			err := runCreate(tc.opts, []string{
				"resource.google_container_node_pool.nodes1.node_config.disk_size_gb",
				"100",
				filename,
			})
			if err != nil {
				t.Fatal(err)
			}

			diff := cmp.Diff(tc.want, readFile(t, filename), cmpopts.AcyclicTransformer("multiline", func(s string) []string {
				return strings.Split(s, "\n")
			}))

			if diff != "" {
				t.Fatalf("Create mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
