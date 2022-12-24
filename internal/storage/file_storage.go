package storage

import (
	"encoding/json"

	"github.com/maxsnegir/url-shortener/internal/utils"
)

type FileData struct {
	Key   string
	Value []byte
}

type FileStorage struct {
	filePath   string
	fileWriter *utils.FileWriter
	fileReader *utils.FileReader
	storage    Storage // In-memory storage
}

func (s *FileStorage) Get(key string) ([]byte, error) {
	return s.storage.Get(key)
}

func (s *FileStorage) Set(key string, value []byte) error {
	if err := s.storage.Set(key, value); err != nil {
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
	return s.fileWriter.Write(encodedData)
}

func (s *FileStorage) loadDumpFromFile() error {
	for {
		fileData := &FileData{}
		encodedData, err := s.fileReader.Read()
		if err != nil {
			return err
		}
		if encodedData == nil {
			break
		}
		if err := json.Unmarshal(encodedData, &fileData); err != nil {
			return err
		}
		if err := s.storage.Set(fileData.Key, fileData.Value); err != nil {
			return nil
		}
	}
	return nil
}
func (s *FileStorage) Shutdown() error {
	if err := s.fileWriter.Close(); err != nil {
		return err
	}
	if err := s.fileReader.Close(); err != nil {
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
		filePath:   filePath,
		fileWriter: fileWriter,
		fileReader: fileReader,
		storage:    NewMapStorage(),
	}
	if err := fileStorage.loadDumpFromFile(); err != nil {
		return fileStorage, LoadingDumbDataError{err: err}
	}
	return fileStorage, nil
}
