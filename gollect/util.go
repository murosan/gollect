// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

func trimQuotes(s string) string {
	return s[1 : len(s)-1]
}
