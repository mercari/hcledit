# hcledit command

[![release][release-badge]][release]
[![docker][docker-badge]][docker]

`hcledit` command is a CLI tool to edit HCL configuration, which exposes [`hcledit`](https://pkg.go.dev/go.mercari.io/hcledit) API as command line interface. You can think this as a sample application built by the package.

## Install

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

[release]: https://github.com/mercari/hcledit/releases
[release-badge]: https://img.shields.io/github/v/release/mercari/hcledit?style=for-the-badge&logo=github

[docker]: https://hub.docker.com/r/mercari/hcledit
[docker-badge]: https://img.shields.io/docker/v/mercari/hcledit?label=docker&sort=semver&style=for-the-badge&logo=docker
