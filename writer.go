// Copyright 2020 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gollect

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

type writer struct {
	io.Writer
	buf      bytes.Buffer
	provider writerProvider
	config   *Config
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *writer) writeForeach() error {
	for _, out := range w.config.OutputPaths {
		wr := w.provider.provide(out)
		if _, err := wr.Write(w.buf.Bytes()); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}
	return nil
}

type stdoutWriter struct{ io.Writer }

func (w *stdoutWriter) Write(p []byte) (int, error) {
	f := os.Stdout
	defer closer(f)
	return f.Write(p)
}

type clipboardWriter struct{ io.Writer }

func (w *clipboardWriter) Write(p []byte) (int, error) {
	if clipboard.Unsupported {
		return 0, errors.New("no support for clipboard")
	}
	return len(p), clipboard.WriteAll(string(p))
}

type fileWriter struct {
	io.Writer
	path string
}

func (w *fileWriter) Write(p []byte) (int, error) {
	file, err := os.Create(w.path)
	if err != nil {
		return 0, err
	}
	defer closer(file)

	return file.WriteAt(p, 0)
}

// io.Writer provider
type writerProvider interface{ provide(s string) io.Writer }

type writerProviderImpl struct{}

func (p *writerProviderImpl) provide(s string) io.Writer {
	switch strings.ToLower(s) {
	case "stdout":
		return &stdoutWriter{}
	case "clipboard":
		return &clipboardWriter{}
	default:
		return &fileWriter{path: s}
	}
}

func closer(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}
