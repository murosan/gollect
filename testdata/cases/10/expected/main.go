package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println(MyReader{Reader: os.Stdin})
}

type MyReader struct{ Reader io.Reader }
