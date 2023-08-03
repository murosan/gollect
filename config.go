// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"gopkg.in/yaml.v3"
)

// Config is a configuration.
type Config struct {
	// path to main.go or glob of main package files
	InputFile string `yaml:"inputFile"`

	// list of output paths
	// filepath, 'stdout' or 'clipboard' are available
	OutputPaths []string `yaml:"outputPaths"`

	// package path prefixes treat as same as builtin packages.
	ThirdPartyPackagePathPrefixes []string `yaml:"thirdPartyPackagePathPrefixes"`

	output io.Writer // used by test
}

func DefaultConfig() *Config {
	return &Config{
		InputFile:   "main.go",
		OutputPaths: []string{"stdout"},
		ThirdPartyPackagePathPrefixes: []string{
			"golang.org/x/exp",
			"github.com/emirpasic/gods",
			"github.com/liyue201/gostl",
			"gonum.org/v1/gonum",
		},
	}
}

// LoadConfig loads config from yaml file.
func LoadConfig(path string) *Config {
	if path == "" {
		return DefaultConfig()
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return UnmarshalConfig(bytes)
}

func UnmarshalConfig(b []byte) *Config {
	c := *DefaultConfig()
	if err := yaml.Unmarshal(b, &c); err != nil {
		panic(err)
	}
	return &c
}

// Validate validates configuration.
func (c *Config) Validate() error {
	if c.InputFile == "" {
		return errors.New("input file is empty")
	}

	for _, out := range c.OutputPaths {
		if strings.ToLower(out) == "clipboard" && clipboard.Unsupported {
			return errors.New("no clipboard option provided for your operating system")
		}
	}

	return nil
}
