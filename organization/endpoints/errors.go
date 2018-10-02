package endpoints

import "net/http"

// Endpoint errors
var (
	ErrorMissingID = errMissingID{}
)

type errMissingID struct{}

func (errMissingID) Error() string {
	return "MISSING_ID"
}

func (errMissingID) StatusCode() int {
	return http.StatusUnauthorized
}
