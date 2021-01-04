# hcledit

`hcledit` is a Go package to edit HCL configurations. Basically, this is just a wrapper of [`hclwrite`](https://pkg.go.dev/github.com/hashicorp/hcl/v2/hclwrite) package which provides low-level features of generating HCL configurations. But `hcledit` allows you to access HCL attribute or block by [`jq`](https://github.com/stedolan/jq)-like query and do various manipulations. See examples of how it works. 

## Install

Use go get:

```bash
$ go get -u github.com/mercari/hcledit
```

## Examples

The following is an HCL configuration which we want to manipulate. 

```hcl
resource "google_container_node_pool" "nodes1" {
   name = "nodes1"

   node_config {
     preemptible  = false
     machine_type = "e2-medium"
   }
}
```

To create a new attribute,

```go
editor, _ := hcledit.Read(filename)
editor.Create("resource.google_container_node_pool.*.node_config.disk_size_gb", "200")
editor.OverWriteFile()
```

To update the existing attribute,

```go
editor, _ := hcledit.Read(filename)
editor.Update("resource.google_container_node_pool.*.node_config.machine_type", "COS")
editor.OverWriteFile()
```

To delete the existing attribute,

```go
editor, _ := hcledit.Read(filename)
editor.Delete("resource.google_container_node_pool.*.node_config.preemptible")
editor.OverWriteFile()
```
