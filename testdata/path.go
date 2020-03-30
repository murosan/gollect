// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"os"
	"path/filepath"
)

var (
	cwd = func() string {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return cwd
	}()

	j    = filepath.Join
	base = j(cwd, "testdata", "codes")

	// FilePaths a set of paths used in test.
	FilePaths = struct {
		Parse,

		A,
		B,
		Pkg1,
		Pkg2,
		Write1,
		Write2 string
	}{
		Parse: j(base, "parse", "main.go"),

		A:      j(base, "a", "main.go"),
		B:      j(base, "b", "*.go"),
		Pkg1:   j(base, "pkg1", "*.go"),
		Pkg2:   j(base, "pkg2", "*.go"),
		Write1: j(base, "writeone", "*.go"),
		Write2: j(base, "writetwo", "*.go"),
	}

	pkgBase = "github.com/murosan/gollect/testdata/codes"

	// PackagePaths a set of package paths used in test.
	PackagePaths = struct {
		Parse1,
		Parse2,

		Pkg1,
		Pkg2 string
	}{
		Parse1: pkgBase + "/parse/apkg",
		Parse2: pkgBase + "/parse/bpkg",

		Pkg1: pkgBase + "/pkg1",
		Pkg2: pkgBase + "/pkg2",
	}
)
