// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"testing"
)

func shouldPanic(t *testing.T, f func(), onfail string) {
	t.Helper()
	defer func() {
		if err := recover(); err == nil {
			// not recovered
			t.Errorf(onfail)
		}
	}()

	f()
}

func shouldNotPanic(t *testing.T, f func(), onfail string) {
	t.Helper()
	defer func() {
		if err := recover(); err != nil {
			// recovered
			t.Errorf(onfail)
		}
	}()

	f()
}
