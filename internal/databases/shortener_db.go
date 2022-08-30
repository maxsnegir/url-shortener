package databases

import (
	"time"
)

type KeyValueDB interface {
	Set(key string, value interface{}, expires time.Duration) error
	Get(key string) (interface{}, error)
}

type URLDateBase map[string]interface{}

func (db URLDateBase) Set(key string, value interface{}, expires time.Duration) error {
	db[key] = value
	return nil
}

func (db URLDateBase) Get(key string) (interface{}, error) {
	value, ok := db[key]
	if !ok {
		return nil, KeyError
	}
	return value, nil
}

func NewURLDateBase() URLDateBase {
	return make(URLDateBase)
}
