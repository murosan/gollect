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
