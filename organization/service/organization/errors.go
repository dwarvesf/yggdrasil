package organization

import "net/http"

// Organization errors
var (
	ErrNameEmpty = errNameEmpty{}
	ErrNotFound  = errNotFound{}
)

type errNotFound struct{}

func (errNotFound) Error() string {
	return "NOT_FOUND"
}

func (errNotFound) StatusCode() int {
	return http.StatusNotFound
}

type errNameEmpty struct{}

func (errNameEmpty) Error() string {
	return "NAME_EMPTY"
}

func (errNameEmpty) StatusCode() int {
	return http.StatusBadRequest
}
