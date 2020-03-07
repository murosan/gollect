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
	base = filepath.Join(cwd, "testdata", "codes")

	// FilePaths a set of paths used in test.
	FilePaths = struct {
		A,
		B,
		Pkg1,
		Pkg2,
		Write1,
		Write2 string
	}{
		A:      filepath.Join(base, "a", "main.go"),
		B:      filepath.Join(base, "b", "*.go"),
		Pkg1:   filepath.Join(base, "pkg1", "*.go"),
		Pkg2:   filepath.Join(base, "pkg2", "*.go"),
		Write1: filepath.Join(base, "write", "*.go"),
		Write2: filepath.Join(base, "writetwo", "*.go"),
	}

	pkgBase = "github.com/murosan/gollect/testdata/codes"

	// PackagePaths a set of package paths used in test.
	PackagePaths = struct {
		Pkg1, Pkg2 string
	}{
		Pkg1: pkgBase + "/pkg1",
		Pkg2: pkgBase + "/pkg2",
	}
)
