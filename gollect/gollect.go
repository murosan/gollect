package gollect

func Main(config *Config) {
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
		panic(err)
	}
	if err := w.writeForeach(); err != nil {
		panic(err)
	}
}
