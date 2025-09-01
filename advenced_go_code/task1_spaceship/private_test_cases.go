package main

import (
	"bytes"
	"errors"
	"io"
)

var (
	closeErr = errors.New("close error")
)

var privateTestCases = []TestCase{
	{
		name: "Большое количество данных",
		check: func(folder Folder, processor Processor) bool {
			batchSize := int64(float64(testMaxFileSize) * maxCompressionRate)
			data := makeData(100 * batchSize)
			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			files, _ := folder.ListFiles()
			if len(files) < 100 {
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
		name: "Чтение исходных данных осуществляется порциями",
		check: func(folder Folder, processor Processor) bool {
			batchSize := int64(float64(testMaxFileSize) * maxCompressionRate)
			data := makeData(20 * batchSize)

			rawData := &mockReader{
				ReadCloser: makeRaw(data),
			}

			err := SaveBackup(rawData, processor, folder)
			if err != nil {
				return false
			}

			return rawData.count >= 20
		},
	},
	{
		name: "Ошибка закрытия стрима для чтения",
		check: func(folder Folder, processor Processor) bool {
			data := makeData(bigDataSize)

			rawData := &errReadCloser{
				ReadCloser: makeRaw(data),
				closeErr:   closeErr,
			}

			err := SaveBackup(rawData, processor, folder)
			if err == nil {
				return false
			}

			return errors.Is(err, closeErr)
		},
	},
	{
		name: "Ошибка закрытия стрима для записи",
		check: func(folder Folder, processor Processor) bool {
			data := makeData(bigDataSize)

			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			res := &errWriteCloser{nopCloser{bytes.NewBuffer(nil)}}
			err = RestoreBackup(res, processor, folder)
			if err == nil {
				return false
			}

			return errors.Is(err, closeErr)
		},
	},
	{
		name: "Ошибка чтения данных из стрима процессора",
		check: func(folder Folder, processor Processor) bool {
			data := makeData(bigDataSize)

			err := SaveBackup(makeRaw(data), processor, folder)
			if err != nil {
				return false
			}

			processor = mocErrProcessor{processor}
			err = RestoreBackup(nopCloser{bytes.NewBuffer(nil)}, processor, folder)
			if err == nil {
				return false
			}

			return true
		},
	},
}

type mockReader struct {
	io.ReadCloser
	count int
}

func (m *mockReader) Read(p []byte) (int, error) {
	m.count++
	return m.ReadCloser.Read(p)
}

type errReadCloser struct {
	io.ReadCloser
	readErr  error
	closeErr error
}

func (r errReadCloser) Close() error { return r.closeErr }

func (r errReadCloser) Read(p []byte) (int, error) {
	if r.readErr != nil {
		return 0, r.readErr
	}

	return r.ReadCloser.Read(p)
}

type errWriteCloser struct {
	io.WriteCloser
}

func (errWriteCloser) Close() error { return closeErr }

type mocErrProcessor struct {
	Processor
}

func (m mocErrProcessor) DecryptAndUncompress(_ io.Reader) io.Reader {
	return errReadCloser{
		readErr: errors.New("read error"),
	}
}
