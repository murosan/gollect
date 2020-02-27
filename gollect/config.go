// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"errors"
	"io/ioutil"
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
}

// LoadConfig loads config from yaml file.
func LoadConfig(path string) *Config {
	var c Config
	if path == "" {
		return &c
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(bytes, &c); err != nil {
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
