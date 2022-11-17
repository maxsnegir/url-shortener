package storage

import (
	"context"
	"encoding/json"

	"github.com/maxsnegir/url-shortener/internal/utils"
)

type FileData struct {
	Key   string
	Value []byte
}

type FileStorage struct {
	FilePath   string
	FileWriter *utils.FileWriter
	FileReader *utils.FileReader
	Storage    Storage // In-memory storage
}

func (s *FileStorage) Get(key string) ([]byte, error) {
	return s.Storage.Get(key)
}

func (s *FileStorage) Set(key string, value []byte) error {
	if err := s.Storage.Set(key, value); err != nil {
		return nil
	}
	fileData := &FileData{
		Key:   key,
		Value: value,
	}
	encodedData, err := json.Marshal(fileData)
	if err != nil {
		return err
	}
	return s.FileWriter.Write(encodedData)
}

func (s *FileStorage) loadDumpFromFile() error {
	for {
		fileData := &FileData{}
		encodedData, err := s.FileReader.Read()
		if err != nil {
			return err
		}
		if encodedData == nil {
			break
		}
		if err := json.Unmarshal(encodedData, &fileData); err != nil {
			return err
		}
		if err := s.Storage.Set(fileData.Key, fileData.Value); err != nil {
			return nil
		}
	}
	return nil
}
func (s *FileStorage) Shutdown(ctx context.Context) error {
	if err := s.FileWriter.Close(); err != nil {
		return err
	}
	if err := s.FileReader.Close(); err != nil {
		return err
	}
	return nil
}

func NewURLFileStorage(filePath string) (Storage, error) {
	fileWriter, err := utils.NewFileWriter(filePath)
	if err != nil {
		return nil, err
	}
	fileReader, err := utils.NewFileReader(filePath)
	if err != nil {
		return nil, err
	}
	fileStorage := &FileStorage{
		FilePath:   filePath,
		FileWriter: fileWriter,
		FileReader: fileReader,
		Storage:    NewMapStorage(),
	}
	if err := fileStorage.loadDumpFromFile(); err != nil {
		return fileStorage, LoadingDumbDataError{err: err}
	}
	return fileStorage, nil
}
