package storage

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maxsnegir/url-shortener/internal/utils"
)

func TestCreatingFileStorage(t *testing.T) {
	filePath := "temp"
	defer func() {
		if err := utils.RemoveFile(filePath); err != nil {
			t.Error(err)
		}
	}()
	_, err := NewURLFileStorage(filePath, NewURLMapStorage())
	t.Run("MapStorage created", func(t *testing.T) {
		require.NoError(t, err, "Error while opening file")
	})
	t.Run("File created and correct", func(t *testing.T) {
		file, err := os.Open(filePath)
		require.NoError(t, err, "Error while opening file")
		require.NotNil(t, file, "File not found")
		require.Equal(t, filePath, file.Name(), "Wrong file name")
	})
}

func TestFileStorageSetData(t *testing.T) {
	filePath := "temp"
	defer func() {
		if err := utils.RemoveFile(filePath); err != nil {
			t.Error(err)
		}
	}()
	storage, err := NewURLFileStorage(filePath, NewURLMapStorage())
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "Full Data",
			key:   "Key 1",
			value: "Value 1",
		},
		{
			name:  "Empty Key",
			key:   "",
			value: "value 2",
		},
		{
			name:  "Empty Value",
			key:   "Key 2",
			value: "",
		},
	}
	t.Run("MapStorage created", func(t *testing.T) {
		require.NoError(t, err, "Error while creating storage")
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Set(tt.key, []byte(tt.value))
			require.NoError(t, err, "Error while setting storage")
			value, err := storage.Get(tt.key)
			require.NoError(t, err, "Error while getting from storage")
			require.Equal(t, tt.value, string(value), "Wrong value from Get() ")
		})
	}
}

func TestStorageIsPersistent(t *testing.T) {
	filePath := "temp"
	defer func() {
		if err := utils.RemoveFile(filePath); err != nil {
			t.Error(err)
		}
	}()
	tests := []UserURL{
		{
			UserID: "12",
			URLData: URLData{
				ShortURL:    "http://localhost:8080/nzsq/",
				OriginalURL: "http://localhost:8080",
			},
		},
		{
			UserID: "10",
			URLData: URLData{
				ShortURL:    "http://localhost:8080/nzsqS1/",
				OriginalURL: "http://practicum.com",
			},
		},
		{
			URLData: URLData{
				OriginalURL: "http://practicum.com",
			},
		},
	}
	firstStorage, err := NewURLFileStorage(filePath, NewURLMapStorage())
	require.NoError(t, err, "Error creating storage")
	for _, dt := range tests {
		require.NoError(t, firstStorage.SetURLData(dt.UserID, dt.URLData), "Error while setting data")
	}

	secondStorage, err := NewURLFileStorage(filePath, NewURLMapStorage())
	require.NoError(t, err, "Error while creating storage")

	for _, tt := range tests {
		t.Run("Data in another storage exist", func(t *testing.T) {
			bytesFromFirstStorage, err := secondStorage.Get(tt.UserID)
			require.NoError(t, err, "error while getting value from first storage")
			bytesFromSecondStorage, err := secondStorage.Get(tt.UserID)
			require.NoError(t, err, "error while getting value from second storage")
			isEqual := bytes.Equal(bytesFromFirstStorage, bytesFromSecondStorage)
			require.True(t, isEqual, "Data in storages are different")
		})
	}
}
