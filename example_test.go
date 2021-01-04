package hcledit_test

import (
	"fmt"
	"strings"

	"github.com/mercari/hcledit"
)

func Example() {
	src := `
resource "google_container_node_pool" "nodes1" {
  name = "nodes1"

  node_config {
    preemptible  = false
    machine_type = "e2-medium"
  }

  timeouts {
    create = "30m"
  }
}

resource "google_container_node_pool" "nodes2" {
  name = "nodes2"

  node_config {
    preemptible  = false
    machine_type = "e2-medium"
  }

  timeouts {
    create = "30m"
  }
}

`
	// Read HCL contents.
	editor, _ := hcledit.Read(strings.NewReader(src), "")

	// Create new attribute on the existing block.
	editor.Create("resource.google_container_node_pool.*.node_config.disk_size_gb", "200")

	// Create new block and add some attributes.
	editor.Create("resource.google_container_node_pool.*.master_auth", hcledit.BlockVal())
	editor.Create("resource.google_container_node_pool.*.master_auth.username", "")
	editor.Create("resource.google_container_node_pool.*.master_auth.password", "")

	// Update existing attributes.
	editor.Update("resource.google_container_node_pool.*.node_config.machine_type", "COS")
	editor.Update("resource.google_container_node_pool.*.node_config.preemptible", true)

	// Delete existing attribute and blocks
	editor.Delete("resource.google_container_node_pool.*.timeouts")

	fmt.Printf("%s", editor.Bytes())
	// Output:
	// resource "google_container_node_pool" "nodes1" {
	//   name = "nodes1"
	//
	//   node_config {
	//     preemptible  = true
	//     machine_type = "COS"
	//     disk_size_gb = "200"
	//   }
	//
	//   master_auth {
	//     username = ""
	//     password = ""
	//   }
	// }
	//
	// resource "google_container_node_pool" "nodes2" {
	//   name = "nodes2"
	//
	//   node_config {
	//     preemptible  = true
	//     machine_type = "COS"
	//     disk_size_gb = "200"
	//   }
	//
	//   master_auth {
	//     username = ""
	//     password = ""
	//   }
	// }
}
