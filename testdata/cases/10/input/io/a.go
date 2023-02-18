package io

import (
	"io"
)

type MyReader struct{ Reader io.Reader }
