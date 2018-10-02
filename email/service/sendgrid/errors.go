package sendgrid

import (
	"net/http"
)

//Error declaration
var (
	ErrNameIsRequired  = errNameIsRequired{}
	ErrEmailIsRequired = errEmailIsRequired{}
)

type errNameIsRequired struct{}

func (errNameIsRequired) Error() string {
	return "Name is required!"
}

func (errNameIsRequired) StatusCode() int {
	return http.StatusBadRequest
}

type errEmailIsRequired struct{}

func (errEmailIsRequired) Error() string {
	return "Email is required!"
}

func (errEmailIsRequired) StatusCode() int {
	return http.StatusBadRequest
}
