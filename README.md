# gollect

![CI](https://github.com/murosan/gollect/workflows/CI/badge.svg?branch=master)

A tool for competitive programming in Go.

[README - 日本語](./docs/README_ja.md)

## Feature

- Extract only what is needed from the code written in multiple packages and output the code that can be submitted.
- Output formatted code.

## Install

```sh
go get -u github.com/murosan/gollect/cmd/gollect
```

## Usage

Suppose you have implemented the following `Max` function in the `lib` package:

```go
package lib

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
```

Then implement `main` package.  
The following code reads two numbers and outputs the larger one.

```go
package main

import (
	"fmt"

	"github.com/your-name/repo-name/lib"
)

func main() {
	var a, b int
	fmt.Scan(&a, &b)

	max := lib.Max(a, b)
	fmt.Println(max)
}
```

Execute `gollect`.

```sh
$ gollect -in ./main.go
```

The code you can submit will be outputted as follow:

```go
package main

import "fmt"

func main() {
	var a, b int
	fmt.Scan(&a, &b)

	max := Max(a, b)
	fmt.Println(max)
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
```

## Configuration

You can see CLI options by executing `gollect -help`.

#### YAML configuration file

Example:

```yml
# Path to main package.
inputFile: ./main.go

# A list of output.
# You can specify stdout, clipboard and filepaths.
outputPaths:
  - stdout
  - clipboard
  - ./out/main.go
```

To run with YAML configuration file, execute like:

```sh
gollect -config ./config.yml
```

## How it works?

Does the followings roughly:

1. list all package-level declarations.
2. find out declarations on which each declaration depends.
3. output all declarations together on which main function in main package depends.

Here's the list of package-level declarations:

- var
- const
- type definition
- function
- method

Methods are also considered package-level declarations.  
And, it doesn't matter if it is exported (begin with a upper case character) or not.  
For example, following is all package-level declarations:

```go
var a = 100
var A = 200
const b = 300
const B = 400
type c struct{}
func (c c) do() {}
func (c *c) Do() {}
type C struct{}
func (C) do() {}
func (*C) Do() {}
type d interface{}
type D interface{}
func e() {}
func E() {}
```

Finally, the declaration which is not used from main function will be ignored.  
Also methods are not exception.

But you may want to keep the method. For example, heap.  
See the example for `IntHeap` on the foillowing page.

https://golang.org/pkg/container/heap/

`Len`, `Less` and the other methods may not be used directly (or indirectly) from the main function, but final code will not work without them.  
To keep them, add `// gollect: keep methods` into comments.
This leaves all `IntHeap`'s methods.

```go
// An IntHeap is a min-heap of ints.
// gollect: keep methods
type IntHeap []int
```
