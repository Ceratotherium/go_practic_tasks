package main

import (
	"bytes"
	"io"
	"sync"
)

type MockFolder struct {
	files map[string]io.Reader
	mutex *sync.Mutex
}

func (f *MockFolder) WriteFile(name string) io.WriteCloser {
	writer := bytes.NewBuffer(make([]byte, 0, maxFileSize))

	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.files[name] = writer

	return nopCloser{writer}
}

func (f *MockFolder) ReadFile(name string) io.ReadCloser {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if file, ok := f.files[name]; ok {
		return io.NopCloser(file)
	}

	return io.NopCloser(bytes.NewBuffer(nil))
}

func (f *MockFolder) ListFiles() ([]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	files := make([]string, 0, len(f.files))
	for file := range f.files {
		files = append(files, file)
	}

	return files, nil
}

func NewMockFolder() *MockFolder {
	return &MockFolder{
		files: make(map[string]io.Reader),
		mutex: &sync.Mutex{},
	}
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
