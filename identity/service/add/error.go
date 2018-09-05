package add

import (
	"net/http"
)

// error types
var (
	ErrInvalidParams = errInvalidParams{}
)

type errInvalidParams struct{}

func (errInvalidParams) Error() string {
	return "invalid parameters"
}
func (errInvalidParams) StatusCode() int {
	return http.StatusBadRequest
}
