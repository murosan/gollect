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
	defer f.Close()
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
	file, err := os.OpenFile(w.path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write(p)
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
