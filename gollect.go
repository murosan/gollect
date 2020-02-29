// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

// Main executes whole program.
func Main(config *Config) error {
	if err := config.Validate(); err != nil {
		return err
	}

	p := NewProgram(config.InputFile)

	// parse ast files and check dependencies
	ParseAll(p)
	AnalyzeForeach(p)

	// mark all used declarations
	next := []ExternalDependencySet{{}}
	next[0].Add("main", "main")
	UseAll(p.Packages(), next)

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
