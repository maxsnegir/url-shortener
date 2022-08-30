package server

import (
	"fmt"
)

type MethodNotAllowedError struct {
	Method string
}

func (e MethodNotAllowedError) Error() string {
	return fmt.Sprintf("Method %s not allowed", e.Method)
}

type InternalServerError struct{}

func (e InternalServerError) Error() string {
	return "Internal Server Error"
}

type NotFoundError struct {
	Url string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Url)
}

type RequestParamsError struct{}

func (e RequestParamsError) Error() string {
	return fmt.Sprint("Url in request body is missing")
}
