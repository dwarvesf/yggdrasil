package group

import "net/http"

// Group errors
var (
	ErrNameEmpty      = errNameEmpty{}
	ErrNotFound       = errNotFound{}
	ErrInvalidStatus  = errInvalidStatus{}
	ErrAlreadyJoined  = errAlreadyJoined{}
	ErrHasNotJoined   = errHasNotJoined{}
	ErrOrgNotFound    = errOrgNotFound{}
	ErrGroupNotActive = errGroupNotActive{}
)

type errInvalidStatus struct{}

func (errInvalidStatus) Error() string {
	return "INVALID_STATUS"
}

func (errInvalidStatus) StatusCode() int {
	return http.StatusBadRequest
}

type errAlreadyJoined struct{}

func (errAlreadyJoined) Error() string {
	return "ALREADY_JOINED"
}

func (errAlreadyJoined) StatusCode() int {
	return http.StatusBadRequest
}

type errHasNotJoined struct{}

func (errHasNotJoined) Error() string {
	return "HAS_NOT_JOINED"
}

func (errHasNotJoined) StatusCode() int {
	return http.StatusBadRequest
}

type errGroupNotActive struct{}

func (errGroupNotActive) Error() string {
	return "GROUP_NOT_ACTIVE"
}

func (errGroupNotActive) StatusCode() int {
	return http.StatusBadRequest
}

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
