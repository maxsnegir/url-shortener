package services

import "fmt"

type OriginalUrlNotFound struct {
	UrlId string
}

func (e OriginalUrlNotFound) Error() string {
	return fmt.Sprintf("Requeted url id = '%s' not found", e.UrlId)
}
