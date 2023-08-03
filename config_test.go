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
		InputFile:                     "abc/def/main.go",
		OutputPaths:                   []string{"stdout", "clipboard"},
		ThirdPartyPackagePathPrefixes: DefaultConfig().ThirdPartyPackagePathPrefixes,
	}
	actual := LoadConfig(filepath.Join(base, "valid.yml"))

	if !reflect.DeepEqual(want, actual) {
		t.Errorf("\n[want]\n%v\n[actual]\n%v", want, actual)
	}

	shouldPanic(t, func() {
		LoadConfig(filepath.Join(base, "invalid.yml"))
	}, "should fail")
}

func TestUnmarshalConfig(t *testing.T) {
	cases := []struct {
		in   string
		want *Config
	}{
		{
			in:   ``,
			want: DefaultConfig(),
		},
		{
			in: `inputFile: main.go
outputPaths:
  - tmp.go
thirdPartyPackagePathPrefixes:
  - golang.org/x/exp
`,
			want: &Config{
				InputFile:                     "main.go",
				OutputPaths:                   []string{"tmp.go"},
				ThirdPartyPackagePathPrefixes: []string{"golang.org/x/exp"},
			},
		},
		{
			in: `inputFile: main.go
outputPaths: ["tmp.go"]
thirdPartyPackagePathPrefixes: [golang.org/x/exp]
`,
			want: &Config{
				InputFile:                     "main.go",
				OutputPaths:                   []string{"tmp.go"},
				ThirdPartyPackagePathPrefixes: []string{"golang.org/x/exp"},
			},
		},
		{
			in: `inputFile: main.go
outputPaths: ["tmp.go"]
thirdPartyPackagePathPrefixes: []
`,
			want: &Config{
				InputFile:                     "main.go",
				OutputPaths:                   []string{"tmp.go"},
				ThirdPartyPackagePathPrefixes: []string{},
			},
		},

		{
			in: `inputFile: tmp/main.go`,
			want: &Config{
				InputFile:                     "tmp/main.go",
				OutputPaths:                   DefaultConfig().OutputPaths,
				ThirdPartyPackagePathPrefixes: DefaultConfig().ThirdPartyPackagePathPrefixes,
			},
		},
	}

	for i, c := range cases {
		config := UnmarshalConfig([]byte(c.in))
		if !reflect.DeepEqual(config, c.want) {
			t.Errorf("at:%d\n[want]\n%v\n[actual]\n%v", i, c.want, config)
		}
	}
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
