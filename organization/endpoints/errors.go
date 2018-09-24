package endpoints

import "net/http"

var (
	//ErrorInvalidLogin return login fail message
	ErrorInvalidLogin = errInvalidLogin{}
)

type errInvalidLogin struct{}

func (errInvalidLogin) Error() string {
	return "User name and password is invalid"
}

func (errInvalidLogin) StatusCode() int {
	return http.StatusUnauthorized
}
