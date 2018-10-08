package referral

import (
	"net/http"
)

var (
	//ErrFromUserIDEmpty ..
	ErrFromUserIDEmpty = errFromUserIDEmpty{}
	//ErrToEmailEmpty ..
	ErrToEmailEmpty = errToEmailEmpty{}
	//ErrCodeEmpty ..
	ErrCodeEmpty = errCodeEmpty{}
	//ErrEmailFormat ..
	ErrEmailFormat = errEmailFormat{}
)

type errFromUserIDEmpty struct{}

func (errFromUserIDEmpty) Error() string {
	return "From User ID empty"
}

func (errFromUserIDEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errToEmailEmpty struct{}

func (errToEmailEmpty) Error() string {
	return "Email empty"
}

func (errToEmailEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errCodeEmpty struct{}

func (errCodeEmpty) Error() string {
	return "Code empty"
}

func (errCodeEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errEmailFormat struct{}

func (errEmailFormat) Error() string {
	return "Email format error"
}

func (errEmailFormat) StatusCode() int {
	return http.StatusBadRequest
}
