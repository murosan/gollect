// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"go/token"
)

// Program is a container of information that is necessary across packages.
type Program struct {
	fset *token.FileSet
	iset *ImportSet
	dset *DeclSet
	pset PackageSet
}

// NewProgram returns new Program.
func NewProgram() *Program {
	return &Program{
		fset: token.NewFileSet(),
		iset: NewImportSet(),
		dset: NewDeclSet(),
		pset: make(PackageSet),
	}
}

// FileSet returns fileset.
func (p *Program) FileSet() *token.FileSet { return p.fset }

// ImportSet returns import set.
func (p *Program) ImportSet() *ImportSet { return p.iset }

// DeclSet returns declaration set.
func (p *Program) DeclSet() *DeclSet { return p.dset }

// PackageSet returns packages.
func (p *Program) PackageSet() PackageSet { return p.pset }
