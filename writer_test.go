// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"testing"
)

func TestWriterProvider(t *testing.T) {
	var provider writerProvider = &writerProviderImpl{}

	for i, s := range []string{
		"clipboard",
		"CLIPBOARD",
		"Clipboard",
		"ClipBoard",
	} {
		switch v := provider.provide(s).(type) {
		case *clipboardWriter:
		// ok
		default:
			t.Errorf("at: %d, want: %s, actual: %T", i, "clipboardWriter", v)
		}
	}

	for i, s := range []string{
		"stdout",
		"STDOUT",
		"Stdout",
		"StdOut",
	} {
		switch v := provider.provide(s).(type) {
		case *stdoutWriter:
		// ok
		default:
			t.Errorf("at: %d, want: %s, actual: %T", i, "stdoutWriter", v)
		}
	}

	for i, s := range []string{"file", "abc"} {
		switch v := provider.provide(s).(type) {
		case *fileWriter:
		// ok
		default:
			t.Errorf("at: %d, want: %s, actual: %T", i, "fileWriter", v)
		}
	}
}
