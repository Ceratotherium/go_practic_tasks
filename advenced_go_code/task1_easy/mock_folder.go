package main

import (
	"bytes"
	"io"
)

type MockFolder struct {
	files map[string]io.Reader
}

func (f MockFolder) WriteFile(name string) io.WriteCloser {
	writer := bytes.NewBuffer(make([]byte, 0, maxFileSize))
	f.files[name] = writer

	return nopCloser{writer}
}

func (f MockFolder) ReadFile(name string) io.ReadCloser {
	if file, ok := f.files[name]; ok {
		return io.NopCloser(file)
	}

	return io.NopCloser(bytes.NewBuffer(nil))
}

func (f MockFolder) ListFiles() ([]string, error) {
	files := make([]string, 0, len(f.files))
	for file := range f.files {
		files = append(files, file)
	}

	return files, nil
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
