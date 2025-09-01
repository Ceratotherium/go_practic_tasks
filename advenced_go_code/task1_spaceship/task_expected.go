package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	minCompressionRate = 1.5 // Минимальный коэффициент сжатия
	maxCompressionRate = 3.0 // Ммаксимальный коэффициент сжатия
	maxProc            = 3
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
	defer func() {
		err := raw.Close() // Зависит от реализации (но обычно ошибку скипают из-за пустышки)
		if resErr == nil {
			resErr = err
		}
	}()

	getBatch := readBatch(raw)

	errGroup, ctx := errgroup.WithContext(context.Background())
	for range maxProc {
		errGroup.Go(func() error {
			return saveByBatch(ctx, getBatch, processor, folder)
		})
	}

	return errGroup.Wait()
}

func saveByBatch(
	ctx context.Context,
	getBatch func(*bytes.Buffer, int64) (int64, error),
	processor Processor,
	folder Folder,
) error {
	batchSize := int64(float64(maxFileSize) * maxCompressionRate)
	buffer := bytes.NewBuffer(make([]byte, 0, batchSize))

	var (
		err          error
		readErr      error
		batchNumber  int64
		writeErrorCh chan error
	)
	getWriteError := func() error {
		if writeErrorCh == nil {
			return nil
		}

		return <-writeErrorCh
	}

	for ctx.Err() == nil && !errors.Is(readErr, io.EOF) {
		batchNumber, readErr = getBatch(buffer, batchSize)
		if readErr != nil && !errors.Is(readErr, io.EOF) || buffer.Len() == 0 {
			break
		}

		processedData := processor.CompressAndEncrypt(buffer)
		if err = getWriteError(); err != nil {
			return err
		}

		if writeErrorCh, err = writeBatch(processedData, folder, batchNumber); err != nil {
			return err
		}
	}

	if err = getWriteError(); err != nil {
		return err
	}

	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return readErr
	}

	return nil
}

func readBatch(raw io.Reader) func(*bytes.Buffer, int64) (int64, error) {
	var batchNumber int64
	mutex := &sync.Mutex{}

	return func(buffer *bytes.Buffer, size int64) (int64, error) {
		buffer.Reset()

		mutex.Lock()
		defer mutex.Unlock()

		n, err := io.CopyN(buffer, raw, size)
		if n == 0 {
			return 0, io.EOF
		}

		currentBatch := batchNumber
		batchNumber++

		return currentBatch, err
	}
}

func writeBatch(data io.Reader, folder Folder, batchNumber int64) (chan error, error) {
	saveData, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	res := make(chan error, 1)
	errGroup := &errgroup.Group{}

	part := int64(0)
	for dataStart := 0; dataStart < len(saveData); dataStart += maxFileSize {
		fileName := makeFileName(batchNumber, part)
		part++

		errGroup.Go(func() error {
			return writeFile(saveData[dataStart:min(dataStart+maxFileSize, len(saveData))], folder, fileName)
		})
	}

	go func() {
		defer close(res)
		res <- errGroup.Wait()
	}()

	return res, nil
}

func writeFile(data []byte, folder Folder, fileName string) (err error) {
	remoteFile := folder.WriteFile(fileName)

	defer func() {
		err = remoteFile.Close() // Зависит от реализации (но обычно ошибку скипают из-за пустышки)
	}()

	if _, err = remoteFile.Write(data); err != nil {
		return err
	}

	return nil
}

func RestoreBackup(
	raw io.WriteCloser,
	processor Processor,
	folder Folder,
) (resErr error) {
	defer func() {
		err := raw.Close()
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

	filesByParts, err := reorderFiles(files)
	if err != nil {
		return err
	}

	processedDataCh := make([]chan io.Reader, 0, maxProc)

	errGroup, ctx := errgroup.WithContext(context.Background())
	for treadIndex := range maxProc {
		processedDataCh = append(processedDataCh, make(chan io.Reader, 1))

		errGroup.Go(func() error {
			defer close(processedDataCh[treadIndex])

			var procErr error
			loadGroup := &sync.WaitGroup{}
			procGroup := &sync.WaitGroup{}

			buffer := bytes.NewBuffer(make([]byte, 0, maxFileSize*maxCompressionRate))

			loadedFiles := make([]io.ReadCloser, 0, maxPartsCount())

			for fileIndex := treadIndex; fileIndex < len(filesByParts); fileIndex += maxProc {
				if ctx.Err() != nil {
					break
				}

				loadedFiles = loadedFiles[:len(filesByParts[fileIndex])]

				for index, part := range filesByParts[fileIndex] {
					wgGo(loadGroup, func() {
						loadedFiles[index] = folder.ReadFile(part)
					})
				}

				loadGroup.Wait()
				procGroup.Wait()

				buffer.Reset()
				for index := range loadedFiles {
					if _, procErr = io.Copy(buffer, loadedFiles[index]); procErr != nil {
						break
					}
				}

				wgGo(procGroup, func() {
					processedDataCh[treadIndex] <- processor.DecryptAndUncompress(buffer)
				})
			}

			procGroup.Wait()

			return procErr
		})
	}

	errGroup.Go(func() error {
		if err := restoreFromParts(raw, processedDataCh); err != nil {
			// Если при восстановлении произошла ошибка, мы не можем просто забить на это и свернуться.
			// Нам нужно добчитать данные из канала, иначе может возникнуть дедлок
			errGroup.Go(func() error {
				for _, ch := range processedDataCh {
					for range ch {
					}
				}
				return nil
			})

			return err
		}

		return nil
	})

	return errGroup.Wait()
}

func restoreFromParts(raw io.Writer, processedDataCh []chan io.Reader) (err error) {
	for {
		closed := 0

		for treadIndex := range maxProc {
			stream, ok := <-processedDataCh[treadIndex]
			if !ok {
				closed++
				continue
			}

			n, err := io.Copy(raw, stream)
			if err != nil {
				return err
			}
			if n == 0 {
				return errors.New("no data written")
			}
		}

		if closed == maxProc {
			break
		}
	}
	return nil
}

func reorderFiles(files []string) ([][]string, error) {
	existFiles := make(map[string]struct{}, len(files))
	result := make([][]string, 0, len(files))

	for _, file := range files {
		existFiles[file] = struct{}{}
	}

	maxPartCount := maxPartsCount()
	totalCount := 0

	for batchNumber := 0; batchNumber < len(files) && totalCount < len(files); batchNumber++ {
		parts := make([]string, 0, maxPartCount)
		for part := range maxPartCount {
			fileName := makeFileName(int64(batchNumber), part)

			_, ok := existFiles[fileName]
			if ok {
				parts = append(parts, fileName)
				totalCount++
			} else if part != maxPartCount-1 {
				return nil, errors.New("part of backup file is skipped")
			}
		}

		if len(parts) == 0 {
			break
		}

		result = append(result, parts)
	}

	if totalCount != len(files) {
		return nil, errors.New("backup batch is skipped")
	}

	return result, nil
}

func maxPartsCount() int64 {
	return int64(math.Floor(maxCompressionRate / minCompressionRate))
}

func makeFileName(butchNumber int64, part int64) string {
	return fmt.Sprintf("backup_%d", (butchNumber+1)*10+part+1)
}

func wgGo(wg *sync.WaitGroup, cb func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		cb()
	}()
}
