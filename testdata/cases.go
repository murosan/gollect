// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var (
	join = filepath.Join

	cwd = func() string {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return cwd
	}()

	base     = join(cwd, "testdata", "cases")
	input    = join("input", "main.go")
	expected = join("expected", "main.go")
	actual   = join("actual", "main.go")

	numOfCases = func() int {
		files, err := ioutil.ReadDir(base)
		if err != nil {
			panic(err)
		}
		return len(files)
	}()

	// Cases is a set of test cases.
	Cases = initCases(numOfCases)
)

type testCase struct {
	Input,
	Expected,
	Actual,
	ActualDir string
}

func newCase(n int) testCase {
	s := strconv.Itoa(n)
	p := func(last string) string { return join(base, s, last) }
	return testCase{
		Input:     p(input),
		Expected:  p(expected),
		Actual:    p(actual),
		ActualDir: filepath.Dir(p(actual)),
	}
}

func initCases(n int) []testCase {
	v := make([]testCase, n)
	for i := 0; i < n; i++ {
		v[i] = newCase(i)
	}
	return v
}
