// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/token"
)

// Program is a container of information that is neccessary across packages.
type Program struct {
	fset     *token.FileSet
	iset     *ImportSet
	packages Packages
}

// NewProgram returns new Program.
func NewProgram() *Program {
	return &Program{
		fset:     token.NewFileSet(),
		iset:     NewImportSet(),
		packages: make(Packages),
	}
}

// FileSet returns fileset.
func (p *Program) FileSet() *token.FileSet { return p.fset }

// ImportSet returns import set.
func (p *Program) ImportSet() *ImportSet { return p.iset }

// Packages returns packages.
func (p *Program) Packages() Packages { return p.packages }
