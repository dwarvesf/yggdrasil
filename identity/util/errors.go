package util

import "net/http"

var (
	//ErrLogin return status 401 in login error
	ErrLogin = errLogin{}
	//ErrorUnauthorize token invalid
	ErrorUnauthorize = errorUnauthorize{}
)

type errLogin struct{}

func (errLogin) Error() string {
	return "User name and password is invalid"
}

func (errLogin) StatusCode() int {
	return http.StatusUnauthorized
}

//ErrAuthentication use for middleware in authentication
type errorUnauthorize struct{}

func (errorUnauthorize) Error() string {
	return "Unauthorize"
}

func (errorUnauthorize) StatusCode() int {
	return http.StatusUnauthorized
}
