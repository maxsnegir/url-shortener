package storages

import "fmt"

const KeyError = DBKeyError("Key does not exist")

type DBKeyError string

func (e DBKeyError) Error() string {
	return string(e)
}

type LoadingDumbDataError struct {
	err error
}

func (e LoadingDumbDataError) Error() string {
	return fmt.Sprintf("Error loading data from file %s", e.err)
}
