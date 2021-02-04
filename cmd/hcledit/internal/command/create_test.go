package command

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunCreate(t *testing.T) {
	cases := map[string]struct {
		opts    *CreateOptions
		want    string
		wantErr bool
	}{
		"WithoutOptionWithAfter": {
			opts: &CreateOptions{
				Type: "string",
			},
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
				Type:      "string",
				WithAfter: "preemptible",
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

			if err := runCreate(tc.opts, []string{
				"resource.google_container_node_pool.nodes1.node_config.disk_size_gb",
				"100",
				filename,
			},
			); (err != nil) != tc.wantErr {
				t.Errorf("runCreate() error = %v, wantErr %v", err, tc.wantErr)
			}

			got := readFile(t, filename)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
