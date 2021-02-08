# hcledit

[![workflow-test][workflow-test-badge]][workflow-test]
[![release][release-badge]][release]
[![docker][docker-badge]][docker]
[![pkg.go.dev][pkg.go.dev-badge]][pkg.go.dev]
[![license][license-badge]][license]

`hcledit` is a Go package to edit HCL configurations. Basically, this is just a wrapper of [`hclwrite`](https://pkg.go.dev/github.com/hashicorp/hcl/v2/hclwrite) package which provides low-level features of generating HCL configurations. But `hcledit` allows you to access HCL attribute or block by [`jq`](https://github.com/stedolan/jq)-like query and do various manipulations. See examples of how it works.

## Install

### Go package

`go get`:

```bash
$ go get -u go.mercari.io/hcledit
```

### Binary

Install binaries via [GitHub Releases][release] or below:

Homebrew:

```bash
$ brew tap mercari/hcledit https://github.com/mercari/hcledit
$ brew install hcledit
```

`go get`:

```bash
$ go get -u go.mercari.io/hcledit/cmd/hcledit
```

Docker:

```bash
$ docker run --rm -it mercari/hcledit hcledit
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
### Go package

To create a new attribute,

```go
editor, _ := hcledit.Read(filename)
editor.Create("resource.google_container_node_pool.*.node_config.disk_size_gb", "200")
editor.OverWriteFile()
```

To update the existing attribute,

```go
editor, _ := hcledit.Read(filename)
editor.Update("resource.google_container_node_pool.*.node_config.image_type", "COS")
editor.OverWriteFile()
```

To delete the existing attribute,

```go
editor, _ := hcledit.Read(filename)
editor.Delete("resource.google_container_node_pool.*.node_config.preemptible")
editor.OverWriteFile()
```

### Binary

To read an attribute,

```console
$ hcledit read 'resource.google_container_node_pool.*.node_config.machine_type' /path/to/file.tf
resource.google_container_node_pool.nodes1.node_config.machine_type e2-medium
```

To create a new attribute,

```console
$ hcledit create 'resource.google_container_node_pool.*.node_config.image_type' 'COS' /path/to/file.tf
```

```diff
resource "google_container_node_pool" "nodes1" {
   name = "nodes1"

   node_config {
     preemptible  = false
     machine_type = "e2-medium"
+    image_type   = "COS"
   }
}
```

To update the existing attribute,

```console
$ hcledit update 'resource.google_container_node_pool.*.node_config.machine_type' 'e2-highmem-2' /path/to/file.tf
```

```diff
resource "google_container_node_pool" "nodes1" {
   name = "nodes1"

   node_config {
     preemptible  = false
-    machine_type = "e2-medium"
+    machine_type = "e2-highmem-2"
   }
}
```

To delete the existing attribute,

```console
$ hcledit delete 'resource.google_container_node_pool.*.node_config.machine_type' /path/to/file.tf
```

```diff
resource "google_container_node_pool" "nodes1" {
   name = "nodes1"

   node_config {
     preemptible  = false
-    machine_type = "e2-medium"
   }
}
```

<!-- badge links -->

[workflow-test]: https://github.com/mercari/hcledit/actions?query=workflow%3ATest
[workflow-test-badge]: https://img.shields.io/github/workflow/status/mercari/hcledit/Test?label=Test&style=for-the-badge&logo=github

[release]: https://github.com/mercari/hcledit/releases
[release-badge]: https://img.shields.io/github/v/release/mercari/hcledit?style=for-the-badge&logo=github

[docker]: https://hub.docker.com/r/mercari/hcledit
[docker-badge]: https://img.shields.io/docker/v/mercari/hcledit?label=docker&sort=semver&style=for-the-badge&logo=docker

[pkg.go.dev]: https://pkg.go.dev/go.mercari.io/hcledit
[pkg.go.dev-badge]: http://bit.ly/pkg-go-dev-badge

[license]: LICENSE
[license-badge]: https://img.shields.io/github/license/mercari/hcledit?style=for-the-badge
