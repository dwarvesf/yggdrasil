package group

import "net/http"

// Group errors
var (
	ErrNameEmpty   = errNameEmpty{}
	ErrNotFound    = errNotFound{}
	ErrOrgNotFound = errOrgNotFound{}
)

type errOrgNotFound struct{}

func (errOrgNotFound) Error() string {
	return "ORGANIZATION_NOT_FOUND"
}

func (errOrgNotFound) StatusCode() int {
	return http.StatusNotFound
}

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
