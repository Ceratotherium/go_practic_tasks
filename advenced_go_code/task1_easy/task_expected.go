package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	compressionRate = 1.5
)

var (
	maxFileSize = 100_000_000 // 100MiB
)

type Processor interface {
	CompressAndEncrypt(stream io.Reader) io.Reader
	DecryptAndUncompress(stream io.Reader) io.Reader
}

type Folder interface {
	WriteFile(name string) io.WriteCloser
	ReadFile(name string) io.ReadCloser
	ListFiles() ([]string, error)
}

func SaveBackup(
	raw io.ReadCloser,
	processor Processor,
	folder Folder,
) (resErr error) {
	batchSize := int64(float64(maxFileSize) * compressionRate)

	defer func() {
		err := raw.Close() // Зависит от реализации (но обычно ошибку скипают из-за пустышки)
		if resErr == nil {
			resErr = err
		}
	}()

	buffer := bytes.NewBuffer(make([]byte, batchSize))

	for batchNumber := int64(0); ; batchNumber++ {
		buffer.Reset()
		n, err := io.CopyN(buffer, raw, batchSize)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 && errors.Is(err, io.EOF) {
			return nil
		}

		err = func() (batchErr error) {
			remoteFile := folder.WriteFile(makeFileName(batchNumber))

			defer func() {
				batchErr = remoteFile.Close() // Зависит от реализации (но обычно ошибку скипают из-за пустышки)
			}()

			processedData := processor.CompressAndEncrypt(io.LimitReader(buffer, n))

			n, err = io.Copy(remoteFile, processedData)
			if err != nil {
				return err
			}

			if n == 0 {
				return errors.New("failed to copy data")
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}
}

func RestoreBackup(
	raw io.WriteCloser,
	processor Processor,
	folder Folder,
) (resErr error) {
	defer func() {
		err := raw.Close() // Зависит от реализации (но обычно ошибку скипают из-за пустышки)
		if resErr == nil {
			resErr = err
		}
	}()

	files, err := folder.ListFiles()
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("no backup files found")
	}

	if err = reorderFiles(files); err != nil {
		return err
	}

	for _, file := range files {
		err = func() (batchErr error) {
			remoteFile := folder.ReadFile(file)

			defer func() {
				batchErr = remoteFile.Close()
			}()

			processedData := processor.DecryptAndUncompress(remoteFile)
			n, err := io.Copy(raw, processedData)
			if err != nil {
				return err
			}

			if n == 0 {
				return errors.New("failed to copy data")
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

func reorderFiles(files []string) error {
	maxNumber := int64(len(files)) - 1
	numbers := make(map[int64]struct{}, len(files))
	for _, file := range files {
		parts := strings.Split(file, "_")
		if len(parts) != 2 {
			return errors.New("invalid file name format")
		}

		number, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return err
		}

		if _, ok := numbers[number]; ok {
			return errors.New("duplicate file")
		}

		if number < 0 || number > maxNumber {
			return errors.New("invalid file number")
		}

		numbers[number] = struct{}{}
	}

	for number := range numbers {
		files[number] = makeFileName(number)
	}

	return nil
}

func makeFileName(butchNumber int64) string {
	return fmt.Sprintf("backup_%d", butchNumber)
}
