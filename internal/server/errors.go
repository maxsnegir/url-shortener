package server

import (
	"fmt"
)

const InternalServerError = internalServerError("Internal Server Error")

type MethodNotAllowedError struct {
	Method string
}

func (e MethodNotAllowedError) Error() string {
	return fmt.Sprintf("Method %s not allowed", e.Method)
}

type internalServerError string

func (e internalServerError) Error() string {
	return string(e)
}

type NotFoundError struct {
	URL string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.URL)
}
