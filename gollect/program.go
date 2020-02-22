package gollect

import (
	"fmt"
	"go/token"
	"path/filepath"
)

type Program struct {
	fset     *token.FileSet
	iset     ImportSet
	packages Packages
	glob     string
}

func NewProgram(glob string) *Program {
	return &Program{
		fset:     token.NewFileSet(),
		iset:     make(ImportSet),
		packages: make(Packages),
		glob:     glob,
	}
}

func (p *Program) FileSet() *token.FileSet { return p.fset }
func (p *Program) ImportSet() ImportSet    { return p.iset }
func (p *Program) Packages() Packages      { return p.packages }

func (p *Program) FilePaths() []string {
	paths, err := filepath.Glob(p.glob)
	if err != nil {
		panic(fmt.Errorf("parse glob: %w", err))
	}

	return paths
}
