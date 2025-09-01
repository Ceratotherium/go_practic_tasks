package main

import (
	"bytes"
	"errors"
	"io"
	"math"
	"math/rand/v2"
	"sort"
	"time"
)

const (
	smallDataSize   = 100
	mediumDataSIze  = 250
	bigDataSize     = 10000
	testMaxFileSize = 200
)

var errListFile = errors.New("error listing file")

type TestCase struct {
	name  string
	check func(folder Folder, processor Processor) bool
}

var testCases = []TestCase{
	// Публичные тесткейсы
	{
		name: "Размер данных меньше максимального размера файла",
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
		name: "Размер сжатых данных меньше максимального размера файла",
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
		name: "Размер сжатых данных больше максимального размера файла",
		check: func(folder Folder, processor Processor) bool {
			batchSize := int64(float64(maxFileSize) * maxCompressionRate)
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
		name: "Сохраненные файлы не превышают максимальный размер",
		check: func(folder Folder, processor Processor) bool {
			data := makeData(bigDataSize)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			for _, file := range files {
				data, _ = io.ReadAll(folder.ReadFile(file))
				if len(data) > maxFileSize {
					return false
				}
			}

			return true
		},
	},
	{
		name: "Размер батча для шифрования учитывает коэффициет сжатия",
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
				if size > int(maxFileSize*maxCompressionRate) {
					return false
				}
			}

			return true
		},
	},
	{
		name: "Дешифрование выполняется батчами",
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
				if size > int(maxFileSize*maxCompressionRate) {
					return false
				}
			}

			return true
		},
	},
	{
		name: "Нет файлов в папке",
		check: func(folder Folder, processor Processor) bool {
			result := bytes.NewBuffer(nil)
			err := RestoreBackup(nopCloser{result}, processor, folder)

			return err != nil
		},
	},
	{
		name: "Ошибка получения фалов в папке",
		check: func(folder Folder, processor Processor) bool {
			result := bytes.NewBuffer(nil)
			folder = mockFolder{
				Folder:       folder,
				listingError: errListFile,
			}
			err := RestoreBackup(nopCloser{result}, processor, folder)

			return errors.Is(err, errListFile)
		},
	},
	{
		name: "Один из файлов бекапа был удален",
		check: func(folder Folder, processor Processor) bool {
			data := makeData(bigDataSize)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			sort.Strings(files)

			mockFolder := folder.(*MockFolder)
			delete(mockFolder.files, files[0])

			result := bytes.NewBuffer(nil)
			err = RestoreBackup(nopCloser{result}, processor, folder)

			return err != nil
		},
	},
	{
		name: "Параллельное восстановление бекапа",
		check: func(folder Folder, processor Processor) bool {
			data := makeData(bigDataSize)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			if len(files) <= 1 {
				return false
			}

			delay := time.Millisecond * 200
			folder = mockFolder{
				Folder: folder,
				Delay:  delay,
			}

			result := bytes.NewBuffer(nil)

			started := time.Now()
			err = RestoreBackup(nopCloser{result}, processor, folder)
			if err != nil {
				return false
			}

			elapsed := time.Since(started)

			//Предполагаем, что минимум в 2 потока будет работать
			maxBatchCount := bigDataSize / float64(maxFileSize) * maxCompressionRate
			return elapsed < delay*time.Duration(maxBatchCount/2)
		},
	},
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
	listingError error
	Delay        time.Duration
}

func (m mockFolder) ListFiles() ([]string, error) {
	if m.listingError != nil {
		return nil, m.listingError
	}

	return m.Folder.ListFiles()
}

func (f mockFolder) WriteFile(name string) io.WriteCloser {
	time.Sleep(f.Delay)
	return f.Folder.WriteFile(name)
}

func (f mockFolder) ReadFile(name string) io.ReadCloser {
	time.Sleep(f.Delay)
	return f.Folder.ReadFile(name)
}
