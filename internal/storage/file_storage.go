package storage

import (
	"encoding/json"

	"github.com/maxsnegir/url-shortener/internal/utils"
)

type URLData struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"full_url"`
}

// FileStorage сохраняет ссылки в файл и в память.
// Является оберткой над MapURLDataBase, чтобы упростить поиск ссылок.
type FileStorage struct {
	FilePath   string
	FileWriter *utils.FileWriter
	FileReader *utils.FileReader
	Storage    Storage // In memory storage
}

func (s *FileStorage) Get(key string) (string, error) {
	return s.Storage.Get(key)
}

// Set Помещаем ссылки в хранилище.
func (s *FileStorage) Set(key, value string) error {
	// Проверяем, что ссылки нет в памяти, чтобы не записывать в файл повторяющиеся данные
	if _, err := s.Storage.Get(key); err == nil {
		return nil
	}
	if err := s.Storage.Set(key, value); err != nil {
		return err
	}
	urlData := URLData{ShortURL: key, FullURL: value}
	encodedData, err := json.Marshal(urlData)
	if err != nil {
		return err
	}

	return s.FileWriter.Write(encodedData)
}

// loadURLFromFile Вызывается при инициализации хранилища.
// Читает все ссылки из файла и загружает в In-Memory хранилище.
func (s *FileStorage) loadURLFromFile() error {
	for {
		urlData := &URLData{}
		encodedData, err := s.FileReader.Read()
		if err != nil {
			return err
		}
		if encodedData == nil {
			break
		}
		if err := json.Unmarshal(encodedData, &urlData); err != nil {
			return err
		}
		if err := s.Storage.Set(urlData.ShortURL, urlData.FullURL); err != nil {
			return nil
		}
	}
	return nil

}

func (s *FileStorage) Shutdown() error {
	if err := s.FileWriter.Close(); err != nil {
		return err
	}
	if err := s.FileReader.Close(); err != nil {
		return err
	}
	return nil

}

func NewFileStorage(filePath string) (Storage, error) {
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
		Storage:    NewMapURLDataBase(),
	}

	if err := fileStorage.loadURLFromFile(); err != nil {
		return fileStorage, LoadingDumbDataError{err: err}
	}
	return fileStorage, nil
}
