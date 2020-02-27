package gollect

import (
	"fmt"
	"go/token"
	"path/filepath"
)

// Program is a container of information that is neccessary across packages.
type Program struct {
	fset     *token.FileSet
	iset     ImportSet
	packages Packages
	glob     string
}

// NewProgram returns new Program.
func NewProgram(glob string) *Program {
	return &Program{
		fset:     token.NewFileSet(),
		iset:     make(ImportSet),
		packages: make(Packages),
		glob:     glob,
	}
}

// FileSet returns fileset.
func (p *Program) FileSet() *token.FileSet { return p.fset }

// ImportSet returns import set.
func (p *Program) ImportSet() ImportSet    { return p.iset }

// Packages returns packages.
func (p *Program) Packages() Packages      { return p.packages }

// FilePaths returns filepaths of glob.
func (p *Program) FilePaths() []string {
	paths, err := filepath.Glob(p.glob)
	if err != nil {
		panic(fmt.Errorf("parse glob: %w", err))
	}

	return paths
}
