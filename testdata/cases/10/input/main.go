package main

import (
	"fmt"
	"os"

	"github.com/murosan/gollect/testdata/cases/10/input/io"
)

func main() {
	fmt.Println(io.MyReader{Reader: os.Stdin})
}
