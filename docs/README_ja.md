# gollect

Go で競技プログラミングを行うためのツールです。

## 特徴

- 複数パッケージに書かれたコードから、必要なものだけを抽出し、提出可能なコードを出力します
- フォーマット済みのコードを出力します

## インストール

```sh
go get -u github.com/murosan/gollect/cmd/gollect
```

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
$ gollect -main ./main.go
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
