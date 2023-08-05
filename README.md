# gollect

![CI](https://github.com/murosan/gollect/workflows/CI/badge.svg?branch=master)

[README - 日本語](./docs/README_ja.md)

A tool for competitive programming in Go.

Extract only the codes used from `main` function, apply `gofmt` and output into one file.

## Install

```sh
go install github.com/murosan/gollect/cmd/gollect@latest
```

It parses AST at runtime, so it depends on the GoLang's version installed.  
Please reinstall when you upgrade (or downgrade) GoLang.

## Usage

Suppose you have implemented the following `Min`,`Max` functions in the `lib` package:

```go
package lib

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
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

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func main() {
	var a, b int
	fmt.Scan(&a, &b)

	max := Max(a, b)
	fmt.Println(max)
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
```

Import statements of self managed package (`github.com/your-name/repo-name/lib`) and unused `Min` functions are not included.

The package `golang.org/x/exp/constraints` is configured to leave by default.
The details of settings are described later.

## Configuration

You can write configuration file by YAML syntax.  
To specify configuration file, run gollect with `-config` option.

```sh
$ gollect -config config.yml
```

### Default values

```yml
inputFile: main.go
outputPaths:
  - stdout
thirdPartyPackagePathPrefixes:
  - golang.org/x/exp
  - github.com/emirpasic/gods
  - github.com/liyue201/gostl
  - gonum.org/v1/gonum
```

### Options

You can override default values by specifying each option.  
The dafault values are used if the option is omitted.

#### `inputFile`

| key       | type   | description                                                                                                                     | default |
| --------- | ------ | ------------------------------------------------------------------------------------------------------------------------------- | ------- |
| inputFile | string | The filepath `main` function is written.<br>If there are multiple files of main package, you have to specify all files by glob. | main.go |

example:

```yml
inputFile: main.go
```

```yml
inputFile: ./*.go
```

#### `outputPaths`

| key         | type     | description                                                     | default |
| ----------- | -------- | --------------------------------------------------------------- | ------- |
| outputPaths | []string | outputs.<br>available values: `stdout`,`clipboard`,`<filepath>` | stdout  |

example:

```yml
outputPaths:
  - stdout
  - clipboard
  - out/main.go
```

#### `thirdPartyPackagePathPrefixes`

| key                           | type     | description                                                                                                                                | default                                                                                          |
| ----------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------ |
| thirdPartyPackagePathPrefixes | []string | Package-path prefixes that can be used at judge system.<br>The import statements and package selectors specified here will not be deleted. | golang.org/x/exp<br>github.com/emirpasic/gods<br>github.com/liyue201/gostl<br>gonum.org/v1/gonum |

example:

```yml
thirdPartyPackagePathPrefixes:
  - golang.org/x/exp
  - github.com/emirpasic/gods
```

```yml
thirdPartyPackagePathPrefixes: []
```

## Other Specification

### Struct Methods

Finally, the declaration which is not used from `main` function will be ignored.  
Also methods are not exception.

```go
// input
package main

import "sort"

type S[T ~int | ~string] []T

func (s S[T]) Len() int           { return len(s) }
func (s S[T]) Less(i, j int) bool { return s[i] < s[j] }
func (s S[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func main() {
	var s S[int]
	sort.Sort(&s)
}
```

```go
// output
// !! compile error !!
package main

import "sort"

type S[T ~int | ~string] []T

func main() {
	var s S[int]
	sort.Sort(&s)
}
```

The `Len`, `Less` and `Swap` functions have been removed as they are not used directly from the `main` function, but the final code will not work without them.  
There are two ways to leave them.

#### 1. Embed `Interface` in `Struct` field

```go
// input
package main

import "sort"

type S[T ~int | ~string] struct {
	sort.Interface
	data []T
}

func (s *S[T]) Len() int           { return len(s.data) }
func (s *S[T]) Less(i, j int) bool { return s.data[i] < s.data[j] }
func (s *S[T]) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (*S[T]) Unused()              {} // will be removed

func main() {
	var s S[int]
	sort.Sort(&s)
}
```

```go
// output
package main

import "sort"

type S[T ~int | ~string] struct {
	sort.Interface
	data []T
}

func (s *S[T]) Len() int           { return len(s.data) }
func (s *S[T]) Less(i, j int) bool { return s.data[i] < s.data[j] }
func (s *S[T]) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }

func main() {
	var s S[int]
	sort.Sort(&s)
}
```

#### 2. Keep all methods by comment annotation

Write `// gollect: keep methods` in the Struct comment, and all methods will be left.

```go
// input
package main

import "sort"

// gollect: keep methods
type S[T ~int | ~string] []T

func (s S[T]) Len() int           { return len(s) }
func (s S[T]) Less(i, j int) bool { return s[i] < s[j] }
func (s S[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (S[T]) Unused()              {} // will be left

func main() {
	var s S[int]
	sort.Sort(&s)
}
```

```go
// output
package main

import "sort"

type S[T ~int | ~string] []T

func (s S[T]) Len() int           { return len(s) }
func (s S[T]) Less(i, j int) bool { return s[i] < s[j] }
func (s S[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (S[T]) Unused()              {} // will be left

func main() {
	var s S[int]
	sort.Sort(&s)
}
```

### Unsupported Statements

#### `cgo`

```go
import "C" // cannot use
```

#### `dot import`

```go
package main
import . "fmt" // cannot use
func main() { Println() }
```

#### `blank import`

```go
package pkg
func init() {}
```

```go
package main
import _ "github.com/owner/repo/pkg" // cannot use
func main() {}
```
