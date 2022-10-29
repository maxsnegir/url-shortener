package utils

import (
	"bufio"
	"os"
)

const (
	WriteFileMask      = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	ReadFileMask       = os.O_RDONLY | os.O_CREATE
	AllFilePermissions = 0644
)

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func (fw *FileWriter) Write(data []byte) error {
	if _, err := fw.writer.Write(data); err != nil {
		return err
	}

	if err := fw.writer.WriteByte('\n'); err != nil {
		return err
	}
	return fw.writer.Flush()
}

func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

type FileReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func (fr *FileReader) Read() ([]byte, error) {
	if !fr.scanner.Scan() {
		return nil, fr.scanner.Err()
	}
	data := fr.scanner.Bytes()
	return data, nil
}

func (fr *FileReader) Close() error {
	return fr.file.Close()
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, WriteFileMask, AllFilePermissions)
	if err != nil {
		return nil, err
	}
	return &FileWriter{file: file, writer: bufio.NewWriter(file)}, nil
}

func NewFileReader(filename string) (*FileReader, error) {
	file, err := os.OpenFile(filename, ReadFileMask, AllFilePermissions)
	if err != nil {
		return nil, err
	}
	return &FileReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func RemoveFile(filename string) error {
	if err := os.Remove(filename); err != nil {
		return err
	}
	return nil
}
