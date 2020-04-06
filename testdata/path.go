// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"path/filepath"
)

var (
	j     = filepath.Join
	codes = j(cwd, "testdata", "codes")

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
		Parse: j(codes, "parse", "main.go"),

		A:      j(codes, "a", "main.go"),
		B:      j(codes, "b", "*.go"),
		Pkg1:   j(codes, "pkg1", "*.go"),
		Pkg2:   j(codes, "pkg2", "*.go"),
		Write1: j(codes, "writeone", "*.go"),
		Write2: j(codes, "writetwo", "*.go"),
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
