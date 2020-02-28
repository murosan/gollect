// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	base := filepath.Join(".", "testdata", "config")

	want := &Config{
		InputFile:   "abc/def/main.go",
		OutputPaths: []string{"stdout", "clipboard"},
	}
	actual := LoadConfig(filepath.Join(base, "valid.yml"))

	if !reflect.DeepEqual(want, actual) {
		t.Errorf("\n[want]\n%v\n[actual]\n%v", want, actual)
	}

	shouldPanic(t, func() {
		LoadConfig(filepath.Join(base, "invalid.yml"))
	}, "should fail")
}

func TestConfig_Validate(t *testing.T) {
	c1 := &Config{InputFile: "", OutputPaths: []string{"stdout"}}
	c2 := &Config{InputFile: "main.go", OutputPaths: []string{"stdout"}}

	if err := c1.Validate(); err == nil {
		t.Error("want error for empty InputFile field but got nil")
	}

	if err := c2.Validate(); err != nil {
		t.Errorf("want: nil, actual: %v", err)
	}
}
