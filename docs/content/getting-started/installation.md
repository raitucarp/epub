---
title: "Installation"
description: "How to install and import the EPUB library into your Go project."
---

## Requirements

- **Go**: Version 1.25 or later.
- **Operating Systems**: Linux, macOS, Windows (cross-platform compatible).

## Adding the Package

In your Go module directory, install the library using `go get`:

```bash
go get github.com/raitucarp/epub
```

## Importing into Your Code

Import the package in your Go source files:

```go
import (
    "github.com/raitucarp/epub"
)
```

## Verifying Installation

Create a test program `main.go`:

```go
package main

import (
    "fmt"
    "github.com/raitucarp/epub"
)

func main() {
    w := epub.New("test-id")
    w.Title("Test Book")
    fmt.Println("epub package installed successfully!")
}
```

Run it:

```bash
go run main.go
```
