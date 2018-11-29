package organization

import "net/http"

// Organization errors
var (
	ErrNameEmpty             = errNameEmpty{}
	ErrNotFound              = errNotFound{}
	ErrAlreadyJoined         = errAlreadyJoined{}
	ErrHasNotJoined          = errHasNotJoined{}
	ErrOrganizationNotActive = errOrganizationNotActive{}
)

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

type errOrganizationNotActive struct{}

func (errOrganizationNotActive) Error() string {
	return "ORG_NOT_ACTIVE"
}

func (errOrganizationNotActive) StatusCode() int {
	return http.StatusBadRequest
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
