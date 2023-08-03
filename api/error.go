package api

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrInvalidId() Error {
	return NewError(http.StatusBadRequest, "invalid id given")
}

func ErrResourceNotFound(res string) Error {
	return NewError(http.StatusNotFound, fmt.Sprintf("resource %v not found", res))
}

func ErrUnAuthorized() Error {
	return NewError(http.StatusUnauthorized, "unauthorized request")
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid json request",
	}
}
