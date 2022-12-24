package services

import "fmt"

type URLIsNotValidError struct {
	URL string
}

func (e URLIsNotValidError) Error() string {
	return fmt.Sprintf("URL %s is not valid", e.URL)
}
