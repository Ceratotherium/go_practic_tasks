package main

import (
	"bytes"
	"io"
)

type MockProcessor struct {
	data map[string][]byte
}

func (m MockProcessor) CompressAndEncrypt(stream io.Reader) io.Reader {
	data, _ := io.ReadAll(stream)

	resData := data[:int64(float64(len(data))/compressionRate)]
	m.data[string(resData)] = data

	return bytes.NewReader(m.invert(resData))
}

func (m MockProcessor) DecryptAndUncompress(stream io.Reader) io.Reader {
	data, _ := io.ReadAll(stream)

	return bytes.NewReader(m.data[string(m.invert(data))])
}

func (m MockProcessor) invert(data []byte) []byte {
	result := make([]byte, 0, len(data))
	for index := range data {
		result = append(result, ^data[index])
	}

	return result
}

func NewMockProcessor() *MockProcessor {
	return &MockProcessor{
		data: make(map[string][]byte),
	}
}
