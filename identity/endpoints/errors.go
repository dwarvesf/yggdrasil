package endpoints

import "net/http"

var (
	//ErrorInvalidLogin return login fail message
	ErrorInvalidLogin = errInvalidLogin{}
	//ErrorNotActive return token invalid message
	ErrorNotActive = errNotActive{}
	//ErrUnauthorize ...
	ErrUnauthorize = errUnauthorize{}
)

type errInvalidLogin struct{}

func (errInvalidLogin) Error() string {
	return "User name and password is invalid"
}

func (errInvalidLogin) StatusCode() int {
	return http.StatusUnauthorized
}

type errNotActive struct{}

func (errNotActive) Error() string {
	return "Account hasn't actived yet"
}

func (errNotActive) StatusCode() int {
	return http.StatusUnauthorized
}

type errUnauthorize struct{}

func (errUnauthorize) Error() string {
	return "Unauthorize"
}

func (errUnauthorize) StatusCode() int {
	return http.StatusUnauthorized
}
