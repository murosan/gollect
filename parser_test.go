// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"path/filepath"
	"testing"

	"github.com/murosan/gollect/testdata"
)

func TestParseAll(t *testing.T) {
	program := NewProgram()
	paths, _ := filepath.Glob(testdata.FilePaths.Parse)

	if program.PackageSet() == nil ||
		program.ImportSet() == nil ||
		len(paths) == 0 ||
		len(program.PackageSet()) != 0 {
		t.Fatalf("something is wrong. %v", program)
	}

	ParseAll(program, "main", paths)
	packages := program.PackageSet()

	if len(packages) != 3 {
		t.Errorf("len=%d, packages=%v", len(packages), packages)
	}

	for i, c := range []struct {
		path  string
		files int
	}{
		{path: "main", files: 1},
		{path: testdata.PackagePaths.Parse1, files: 2},
		{path: testdata.PackagePaths.Parse2, files: 1},
	} {
		v, ok := packages[c.path]
		if !ok {
			t.Errorf("key not found. at=%d, key=%s", i, c.path)
		}

		if len(v.files) != c.files {
			t.Errorf("files count. at=%d, want=%d, actual=%d", i, c.files, len(v.files))
		}
	}
}
