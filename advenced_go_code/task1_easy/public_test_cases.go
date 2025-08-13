package main

import (
	"bytes"
	"errors"
	"io"
	"math"
	"math/rand/v2"
	"sort"
)

const (
	smallDataSize   = 100
	mediumDataSIze  = 250
	bigDataSize     = 10000
	testMaxFileSize = 200
)

var errListFile = errors.New("error listing file")

type TestCase struct {
	name   string
	folder MockFolder
	check  func(folder Folder, processor Processor) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name:   "Размер данных меньше максимального размера файла",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {
			data := makeData(smallDataSize)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			if len(files) != 1 {
				return false
			}

			result := bytes.NewBuffer(nil)
			err = RestoreBackup(nopCloser{result}, processor, folder)
			if err != nil {
				return false
			}

			return bytes.Equal(data, result.Bytes())
		},
	},
	{
		name:   "Размер сжатых данных меньше максимального размера файла",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {
			data := makeData(mediumDataSIze)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			if len(files) != 1 {
				return false
			}

			result := bytes.NewBuffer(nil)
			err = RestoreBackup(nopCloser{result}, processor, folder)
			if err != nil {
				return false
			}

			return bytes.Equal(data, result.Bytes())
		},
	},
	{
		name:   "Размер сжатых данных больше максимального размера файла",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {
			batchSize := int64(float64(testMaxFileSize) * compressionRate)
			data := makeData(2*batchSize - 1)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			if len(files) <= 1 {
				return false
			}

			result := bytes.NewBuffer(nil)
			err = RestoreBackup(nopCloser{result}, processor, folder)
			if err != nil {
				return false
			}

			return bytes.Equal(data, result.Bytes())
		},
	},
	// Тесткейсы в помощь
	{
		name:   "Сохраненные файлы не превышают максимальный размер",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {
			data := makeData(bigDataSize)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			for _, file := range files {
				data, _ = io.ReadAll(folder.ReadFile(file))
				if len(data) > testMaxFileSize {
					return false
				}
			}

			return true
		},
	},
	{
		name:   "Размер батча для шифрования учитывает коэффициет сжатия",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {
			testProcessor := &mockProcessor{
				Processor:  processor,
				batchSizes: make([]int, 0),
			}

			data := makeData(bigDataSize)
			err := SaveBackup(makeRaw(data), testProcessor, folder)
			if err != nil {
				return false
			}

			for _, size := range testProcessor.batchSizes {
				if size > int(testMaxFileSize*compressionRate) {
					return false
				}
			}

			return true
		},
	},
	{
		name:   "Дешифрование выполняется батчами",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {

			data := makeData(bigDataSize)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			testProcessor := &mockProcessor{
				Processor:  processor,
				batchSizes: make([]int, 0),
			}

			result := bytes.NewBuffer(nil)
			err = RestoreBackup(nopCloser{result}, testProcessor, folder)

			for _, size := range testProcessor.batchSizes {
				if size > int(testMaxFileSize*compressionRate) {
					return false
				}
			}

			return true
		},
	},
	{
		name:   "Нет файлов в папке",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {
			result := bytes.NewBuffer(nil)
			err := RestoreBackup(nopCloser{result}, processor, folder)

			return err != nil
		},
	},
	{
		name:   "Ошибка получения фалов в папке",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {
			result := bytes.NewBuffer(nil)
			err := RestoreBackup(nopCloser{result}, processor, mockFolder{folder})

			return errors.Is(err, errListFile)
		},
	},
	{
		name:   "Один из файлов бекапа был удален",
		folder: makeFolder(),
		check: func(folder Folder, processor Processor) bool {
			data := makeData(bigDataSize)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			sort.Strings(files)

			mockFolder := folder.(MockFolder)
			delete(mockFolder.files, files[0])

			result := bytes.NewBuffer(nil)
			err = RestoreBackup(nopCloser{result}, processor, folder)

			return err != nil
		},
	},
}

func makeFolder() MockFolder {
	return MockFolder{
		files: make(map[string]io.Reader),
	}
}

func makeRaw(data []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(data))
}

func makeData(size int64) []byte {
	data := make([]byte, 0, size)

	for range size {
		data = append(data, byte(rand.Int()%math.MaxInt8))
	}

	return data
}

type mockProcessor struct {
	Processor
	batchSizes []int
}

func (m *mockProcessor) CompressAndEncrypt(stream io.Reader) io.Reader {
	data, _ := io.ReadAll(stream)
	m.batchSizes = append(m.batchSizes, len(data))
	return m.Processor.CompressAndEncrypt(bytes.NewReader(data))
}

func (m *mockProcessor) DecryptAndUncompress(stream io.Reader) io.Reader {
	data, _ := io.ReadAll(stream)
	m.batchSizes = append(m.batchSizes, len(data))
	return m.Processor.DecryptAndUncompress(bytes.NewReader(data))
}

type mockFolder struct {
	Folder
}

func (m mockFolder) ListFiles() ([]string, error) {
	return nil, errListFile
}
