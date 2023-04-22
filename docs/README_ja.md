# gollect

Go で競技プログラミングを行うためのツールです。

## 特徴

- 複数パッケージに書かれたコードから、必要なものだけを抽出し、提出可能なコードを出力します
- フォーマット済みのコードを出力します

## インストール

```sh
go install github.com/murosan/gollect/cmd/gollect@latest
```

実行時に AST をパースしている関係で、インストール時の Go のバージョンに依存しています。  
Go のバージョンを変更したときは、再インストールしてください。

## 使い方

以下のような `Max` 関数を `lib` パッケージに実装したとします。

```go
package lib

func Max(a, b int) int {
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

## 設定

cli オプションは `gollect -help` で表示できます。

#### YAML 設定ファイル

設定例

```yml
# main パッケージのファイルパス
inputFile: ./main.go

# 出力先のリスト
# 標準出力、クリップボード、ファイルパスを指定できます。
outputPaths:
  - stdout
  - clipboard
  - ./out/main.go
```

YAML の設定ファイルを指定して実行するには以下のようにします。

```sh
gollect -config ./config.yml
```

## 何を行っているか

おおまかに以下のことを行います

1. パッケージレベルの宣言をリストアップする
2. それぞれの宣言が依存している宣言を調べる
3. main パッケージの main 関数が依存している宣言をすべて一つのファイルにまとめて出力する

パッケージレベルの宣言の一覧は以下です。

- var
- const
- type 定義
- 関数
- メソッド

メソッドもパッケージレベルの宣言とみなします。  
また、Exported(大文字で始まる)かどうかは関係ありません。  
例えば以下はすべてパッケージレベルの宣言です。

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

最終的に main 関数から使用されていない宣言は無視されます。  
メソッドも例外ではありません。

しかし、メソッドは残したい場合があると思います。
例えば heap です。  
以下のページの `IntHeap` の例を見てください。

https://golang.org/pkg/container/heap/

`Len` や `Less` などのメソッドは main 関数から直接(または間接的に)使用されることはないかもしれませんが、残さなければコードが機能しません。  
これらを残すには、コメントに `// gollect: keep methods` を追記してください。  
これで `IntHeap` のメソッドはすべて残されます。

```go
// An IntHeap is a min-heap of ints.
// gollect: keep methods
type IntHeap []int
```
