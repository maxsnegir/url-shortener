package storage

import (
	"encoding/json"

	"github.com/maxsnegir/url-shortener/internal/utils"
)

type UserURL struct {
	UserID string `json:"user_id"`
	URLData
}

type URLFileStorage struct {
	FilePath   string
	FileWriter *utils.FileWriter
	FileReader *utils.FileReader
	Storage    URLStorage // In memory storage
}

func (s *URLFileStorage) Get(key string) ([]byte, error) {
	return s.Storage.Get(key)
}

func (s *URLFileStorage) Set(key string, value []byte) error {
	return s.Storage.Set(key, value)
}

func (s *URLFileStorage) GetURLData(userID string) (URLDataList, error) {
	return s.Storage.GetURLData(userID)
}

func (s *URLFileStorage) SetURLData(userID string, urlData URLData) error {
	if err := s.Storage.SetURLData(userID, urlData); err != nil {
		return err
	}
	userWithURL := UserURL{
		UserID:  userID,
		URLData: urlData,
	}
	encodedUserWithURLs, err := json.Marshal(userWithURL)
	if err != nil {
		return err
	}
	return s.FileWriter.Write(encodedUserWithURLs)
}

func (s *URLFileStorage) Shutdown() error {
	if err := s.FileWriter.Close(); err != nil {
		return err
	}
	if err := s.FileReader.Close(); err != nil {
		return err
	}
	return nil
}

func (s *URLFileStorage) loadDumpFromFile() error {
	for {
		userURL := &UserURL{}
		encodedData, err := s.FileReader.Read()
		if err != nil {
			return err
		}
		if encodedData == nil {
			break
		}
		if err := json.Unmarshal(encodedData, &userURL); err != nil {
			return err
		}
		if err := s.Storage.SetURLData(userURL.UserID, userURL.URLData); err != nil {
			return nil
		}
	}
	return nil
}
func configureFileStorage(filePath string, storage URLStorage) (*URLFileStorage, error) {
	fileWriter, err := utils.NewFileWriter(filePath)
	if err != nil {
		return nil, err
	}
	fileReader, err := utils.NewFileReader(filePath)
	if err != nil {
		return nil, err
	}
	fileStorage := &URLFileStorage{
		FilePath:   filePath,
		FileWriter: fileWriter,
		FileReader: fileReader,
		Storage:    storage,
	}
	return fileStorage, nil
}

func NewURLFileStorage(filePath string, mapStorage URLStorage) (URLStorage, error) {
	fileStorage, err := configureFileStorage(filePath, mapStorage)
	if err != nil {
		return nil, err
	}
	if err := fileStorage.loadDumpFromFile(); err != nil {
		return fileStorage, err
	}
	return fileStorage, nil
}
