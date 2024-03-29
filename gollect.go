// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"fmt"
	"path/filepath"
)

// Main executes whole program.
func Main(config *Config) error {
	if err := config.Validate(); err != nil {
		return err
	}

	setThirdPartyPackagePathPrefixes(config.ThirdPartyPackagePathPrefixes)

	p := NewProgram()

	paths, err := filepath.Glob(config.InputFile)
	if err != nil {
		return fmt.Errorf("parse glob: %w", err)
	}

	// parse ast files and check dependencies
	ParseAll(p, "main", paths)
	AnalyzeForeach(p, "main", "main")

	w := &writer{
		config:   config,
		provider: &writerProviderImpl{},
	}

	if err := Write(w, p); err != nil {
		return err
	}
	if err := w.writeForeach(); err != nil {
		return err
	}

	return nil
}
