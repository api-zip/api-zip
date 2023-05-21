# The Zip API Object Framework

[![](https://pkg.go.dev/badge/api.zip.svg)](https://pkg.go.dev/api.zip)
![](https://img.shields.io/static/v1?label=license&message=BSD-3&color=%23385177)
[![Go Report Card](https://goreportcard.com/badge/api.zip)](https://goreportcard.com/report/api.zip)
![Latest release](https://img.shields.io/github/v/release/api-zip/api-zip)

This project contains the low-level object framework Zip API.  It enables
Kubernetes-inspired objects without the hassle of code generation.

## Usage

To define your objects using the framework, simply import `api.zip` and declare
a `zip.Object` along with a Spec and Status to get started.

```go
import "api.zip"

type BookSpec struct {
	Name   string
	Author string
}

type BookState string

const (
	BookStateAvailable   = BookState("available")
	BookStateUnavailable = BookState("unavailable")
)

type BookStatus struct {
	State BookState
}

type (
	Book     = zip.Object[BookSpec, BookStatus]
	BookList = zip.ObjectList[BookSpec, BookStatus]
)
```

## License

The Zip API Object Framework is licensed under [BSD-3-Clause](/LICENSE.md).
