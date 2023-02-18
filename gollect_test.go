// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	dmp "github.com/sergi/go-diff/diffmatchpatch"

	"github.com/murosan/gollect/testdata"
)

func TestGollect(t *testing.T) {
	for i, tc := range testdata.Cases {
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()

			conf := &Config{
				InputFile:   tc.Input,
				OutputPaths: []string{tc.Actual},
			}

			fatal := func(t *testing.T, i int, msg string, err error) {
				t.Helper()
				t.Fatalf("At: %d, %s, %v", i, msg, err)
			}

			if _, err := os.Stat(tc.ActualDir); err != nil {
				if err := os.Mkdir(tc.ActualDir, 0755); err != nil {
					fatal(t, i, "create actual dir", err)
				}
			}

			if err := Main(conf); err != nil {
				fatal(t, i, "call Main", err)
			}

			expected, err := ioutil.ReadFile(tc.Expected)
			if err != nil {
				fatal(t, i, "read expected file", err)
			}

			actual, err := ioutil.ReadFile(tc.Actual)
			if err != nil {
				fatal(t, i, "read actual file", err)
			}

			if !bytes.Equal(expected, actual) {
				diff := dmp.New().DiffMain(string(expected), string(actual), true)
				t.Errorf(`
===================================================================
At: %d
Diff:
%s
===================================================================
`, i, colorDiff(diff))
			}
		})
	}
}

// https://github.com/sergi/go-diff/blob/master/diffmatchpatch/diff.go#L1183
func colorDiff(diffs []dmp.Diff) string {
	var b strings.Builder

	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case dmp.DiffInsert:
			_, _ = b.WriteString("\x1b[32m +")
			_, _ = b.WriteString(text)
			_, _ = b.WriteString("\x1b[0m")
		case dmp.DiffDelete:
			_, _ = b.WriteString("\x1b[31m -")
			_, _ = b.WriteString(text)
			_, _ = b.WriteString("\x1b[0m")
		case dmp.DiffEqual:
			_, _ = b.WriteString(text)
		}
	}

	return b.String()
}
