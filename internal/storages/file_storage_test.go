package storages

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCreatingFileStorage(t *testing.T) {
	filePath := "temp"
	defer func() { _ = os.Remove(filePath) }()
	_, err := NewFileStorage(filePath)
	t.Run("Storage created", func(t *testing.T) {
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
	defer func() { _ = os.Remove(filePath) }()
	storage, err := NewFileStorage(filePath)
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
	t.Run("Storage created", func(t *testing.T) {
		require.NoError(t, err, "Error while creating storage")
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Set(tt.key, tt.value)
			require.NoError(t, err, "Error while setting storage")
			value, err := storage.Get(tt.key)
			require.NoError(t, err, "Error while getting from storage")
			require.Equal(t, tt.value, value, "Wrong value from Get() ")
		})
	}
}

func TestStorageIsPersistent(t *testing.T) {
	filePath := "temp"
	defer func() { _ = os.Remove(filePath) }()
	firstStorage, _ := NewFileStorage(filePath)
	tests := []struct {
		key   string
		value string
	}{
		{
			key:   "Key 1",
			value: "value",
		},
		{
			key:   "Key 2",
			value: "value 2",
		},
		{
			key:   "Key 2",
			value: "value 2",
		},
	}
	for _, dt := range tests {
		_ = firstStorage.Set(dt.key, dt.value)
	}

	secondStorage, _ := NewFileStorage(filePath)

	for _, tt := range tests {
		t.Run("Data in another storage exist", func(t *testing.T) {
			value, err := secondStorage.Get(tt.key)
			require.NoError(t, err, "Error while getting value from second storage.")
			require.Equal(t, tt.value, value, "Data in second storage is wrong.")
		})
	}
}
