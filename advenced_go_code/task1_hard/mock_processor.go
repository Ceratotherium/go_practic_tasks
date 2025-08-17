package main

import (
	"bytes"
	"io"
	"math/rand"
	"sync"
)

type MockProcessor struct {
	data  map[string][]byte
	mutex *sync.Mutex
}

func (m *MockProcessor) CompressAndEncrypt(stream io.Reader) io.Reader {
	data, _ := io.ReadAll(stream)

	currentCompressionRate := minCompressionRate + float64(rand.Int()%((maxCompressionRate-minCompressionRate)*10))/10.0

	resData := data[:int64(float64(len(data))/currentCompressionRate)]

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[string(resData)] = data

	return bytes.NewReader(m.invert(resData))
}

func (m *MockProcessor) DecryptAndUncompress(stream io.Reader) io.Reader {
	data, _ := io.ReadAll(stream)

	m.mutex.Lock()
	defer m.mutex.Unlock()
	return bytes.NewReader(m.data[string(m.invert(data))])
}

func (m *MockProcessor) invert(data []byte) []byte {
	result := make([]byte, 0, len(data))
	for index := range data {
		result = append(result, ^data[index])
	}

	return result
}

func NewMockProcessor() *MockProcessor {
	return &MockProcessor{
		data:  make(map[string][]byte),
		mutex: &sync.Mutex{},
	}
}
