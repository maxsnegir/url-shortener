package services

import "fmt"

type OriginalURLNotFound struct {
	UrlID string
}

func (e OriginalURLNotFound) Error() string {
	return fmt.Sprintf("Requeted url id = '%s' not found", e.UrlID)
}

type URLIsNotValidError struct {
	Url string
}

func (e URLIsNotValidError) Error() string {
	return fmt.Sprintf("URL %s is not valid", e.Url)
}
