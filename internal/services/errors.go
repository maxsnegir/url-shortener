package services

import "fmt"

type OriginalURLNotFound struct {
	URLID string
}

func (e OriginalURLNotFound) Error() string {
	return fmt.Sprintf("Requested url id = '%s' not found", e.URLID)
}

type URLIsNotValidError struct {
	URL string
}

func (e URLIsNotValidError) Error() string {
	return fmt.Sprintf("URL %s is not valid", e.URL)
}
