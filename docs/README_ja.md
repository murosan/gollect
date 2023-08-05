# gollect

Goで競技プログラミングを行うためのツールです。

`main`関数から使用されているコードのみを抽出し、`gofmt`をかけた上で、1 つのファイルとして出力します。

## インストール

```sh
go install github.com/murosan/gollect/cmd/gollect@latest
```

実行時に AST をパースしている関係で、インストール時の Go のバージョンに依存しています。  
Go のバージョンを変更したときは、再インストールしてください。

## 使い方

以下のような `Min`,`Max` 関数を `lib` パッケージに実装したとします。

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

続いて、`main` パッケージを実装します。  
以下は 2 つの数値を読み取って大きい方を出力するコードです。

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

gollect を実行します。

```sh
$ gollect -in ./main.go
```

以下のように、そのまま提出できるコードが出力されます。

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

自作パッケージ(`github.com/your-name/repo-name/lib`)のインポート文や、使用されていない`Min`関数は出力されません。

`golang.org/x/exp/constraints`パッケージはデフォルトで残すように設定されています。
設定内容は後述します。

## 設定

設定ファイルを YAML で書くことができます。  
設定ファイルを cli で指定するには`-config`オプションで指定します。

```sh
$ gollect -config config.yml
```

### デフォルト設定

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

### 設定項目

指定すると、デフォルトの設定を上書きできます。  
オプションが省略された場合、デフォルトの設定が使用されます。

#### `inputFile`

| key       | type   | description                          | default |
| --------- | ------ | ------------------------------------ | ------- |
| inputFile | string | `main`関数があるファイルを指定します | main.go |

example:

```yml
inputFile: main.go
```

#### `outputPaths`

| key         | type     | description                                                               | default |
| ----------- | -------- | ------------------------------------------------------------------------- | ------- |
| outputPaths | []string | 出力先を指定します。<br>設定可能な値: `stdout`,`clipboard`,`ファイルパス` | stdout  |

example:

```yml
outputPaths:
  - stdout
  - clipboard
  - out/main.go
```

#### `thirdPartyPackagePathPrefixes`

| key                           | type     | description                                                                                                                                    | default                                                                                          |
| ----------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| thirdPartyPackagePathPrefixes | []string | ジャッジシステム上で使用できるパッケージのプレフィックスを指定します。<br>ここで指定されたパッケージはインポート文などが削除されずに残ります。 | golang.org/x/exp<br>github.com/emirpasic/gods<br>github.com/liyue201/gostl<br>gonum.org/v1/gonum |

example:

```yml
thirdPartyPackagePathPrefixes:
  - golang.org/x/exp
  - github.com/emirpasic/gods
```

```yml
thirdPartyPackagePathPrefixes: []
```

## その他仕様

### Struct Methods

最終的に main 関数から使用されていない宣言は無視されます。  
メソッドも例外ではありません。

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

`Len`、`Less`、`Swap` メソッドは `main` 関数から直接呼び出されないため削除されていますが、残さなければコードが機能しません。  
これらを残すには、2 つの方法があります。

#### 方法 1. Interface を Struct フィールドに埋め込む

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

#### 方法 2. 全てのメソッドを残す

コメントに `// gollect: keep methods` を書くと、全てのメソッドを残します。

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

### サポートされない動作

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
