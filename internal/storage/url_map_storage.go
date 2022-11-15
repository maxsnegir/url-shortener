package storage

import "encoding/json"

type URLs []URLData

type UserURLs struct {
	UserID string `json:"user_id"`
	URLs   URLs   `json:"urls"`
}

type URLMapStorage struct {
	MapStorage
}

func (s *URLMapStorage) GetURLData(userID string) (URLDataList, error) {
	var userURLData URLDataList
	encodedData, err := s.Get(userID)
	if err != nil {
		if err == KeyError {
			return userURLData, nil
		}
		return userURLData, err
	}
	if err := json.Unmarshal(encodedData, &userURLData); err != nil {
		return userURLData, err
	}
	return userURLData, nil
}

func (s *URLMapStorage) SetURLData(userID string, urlData URLData) error {
	userURLs, err := s.GetURLData(userID)
	if err != nil {
		return nil
	}
	for _, data := range userURLs {
		if data == urlData {
			return nil
		}
	}
	userURLs = append(userURLs, urlData)
	encodedUserURLS, err := json.Marshal(userURLs)
	if err != nil {
		return err
	}
	return s.Set(userID, encodedUserURLS)
}

func NewURLMapStorage() URLStorage {
	return &URLMapStorage{
		MapStorage: configureMapStorage(),
	}
}
