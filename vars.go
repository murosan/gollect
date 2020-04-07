// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

// Annotation is an annotation for doc text.
type Annotation string

func (a Annotation) String() string { return string(a) }

const (
	annotationPrefix = "// gollect: "

	// An annotation for type declaration.
	// If exists in doc comment, all receiver methods will be left,
	// otherwise only the methods called directly will be left.
	keepMethods Annotation = annotationPrefix + "keep methods"
)
