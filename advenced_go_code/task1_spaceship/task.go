//go:build task_template

package main

// Processor предоставляет алгоритм работы
// со сжатием и шифрованием потока
type Processor interface {
	// CompressAndEncrpyt сжимает и шифрует поток
	CompressAndEncrypt(stream io.Reader) io.Reader

	// DecryptAndUncompress расшифровывает и разархивирует поток
	DecryptAndUncompress(stream io.Reader) io.Reader
}

// Folder реализует сетевое хранилище
type Folder interface {
	// WriteFile открывает файл на запись по имени
	WriteFile(name string) io.WriteCloser

	// ReadFile открывает файл на чтение по имени
	ReadFile(name string) io.ReadCloser

	// ListFiles возвращает список файлов в хранилище
	ListFiles() ([]string, error)
}

func SaveBackup(
	raw io.ReadCloser,
	processor Processor,
	folder Folder,
) error {
	// TODO
}

func RestoreBackup(
	raw io.WriteСloser,
	processor Processor,
	folder Folder,
) error {
	// TODO
}
