package storage

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

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

func isDuplicateErr(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation
}

type DuplicateURLErr struct {
	URL string
}

func (e DuplicateURLErr) Error() string {
	return fmt.Sprintf("%s already exists", e.URL)
}

func NewDuplicateError(url string) error {
	return &DuplicateURLErr{URL: url}
}
