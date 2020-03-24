# rewindreader

[![Build Status](https://travis-ci.com/TV4/rewindreader.svg?token=DxcU5hM6x3rPkZHnQyTy&branch=master)](https://travis-ci.com/TV4/rewindreader)
[![Go Report Card](https://goreportcard.com/badge/github.com/TV4/rewindreader)](https://goreportcard.com/report/github.com/TV4/rewindreader)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/TV4/rewindreader)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/TV4/rewindreader#license)

`rewindreader` is a Go package that seamlessly buffers a stream of data to allow
reading part of the stream (for example a header of some kind) and then
rewinding it to be able to read the data in full from the start of the stream.

## Installation
```
go get -u github.com/TV4/rewindreader
```

## Usage
```go
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TV4/rewindreader"
)

func main() {
	rwr := rewindreader.New(strings.NewReader("foo bar baz"))

	io.CopyN(os.Stdout, rwr, 3)
	fmt.Println()

	if err := rwr.Rewind(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	io.Copy(os.Stdout, rwr)
	fmt.Println()
}
```

Output:
```
foo
foo bar baz
```

# License
Copyright (c) 2019-2020 TV4

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
