package user

import "net/http"

var (
	//ErrEmailEmpty check email not empty
	ErrEmailEmpty = errEmailEmpty{}
	//ErrPhoneEmpty check phonenumber not empty
	ErrPhoneEmpty = errPhoneEmpty{}
	//ErrUsernameEmpty check username not empty
	ErrUsernameEmpty = errUsernameEmpty{}
	//ErrEmailFormat check email has struct:xx@xx.xx
	ErrEmailFormat = errEmailFormat{}
	//ErrPhoneNumberFortmat check phone has struct +xx0xxxxxxx
	ErrPhoneNumberFortmat = errPhoneNumberFormat{}
	//ErrUsernameExist ...
	ErrUsernameExist = errUsernameExist{}
	//ErrPhonenumberExist ...
	ErrPhonenumberExist = errPhoneNumberExist{}
	//ErrEmailExist ...
	ErrEmailExist = errEmailExist{}
)

type errIdentity struct{}

func (errIdentity) Error() string {
	return "Identity invalid"
}

func (errIdentity) StatusCode() int {
	return http.StatusBadRequest
}

type errEmailEmpty struct{}

func (errEmailEmpty) Error() string {
	return "Email empty"
}

func (errEmailEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errPhoneEmpty struct{}

func (errPhoneEmpty) Error() string {
	return "Phone empty"
}

func (errPhoneEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errUsernameEmpty struct{}

func (errUsernameEmpty) Error() string {
	return "Username empty"
}

func (errUsernameEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errEmailFormat struct{}

func (errEmailFormat) Error() string {
	return "Email invalid"
}

func (errEmailFormat) StatusCode() int {
	return http.StatusBadRequest
}

type errPhoneNumberFormat struct{}

func (errPhoneNumberFormat) Error() string {
	return "Phonenumber invalid"
}

func (errPhoneNumberFormat) StatusCode() int {
	return http.StatusBadRequest
}

type errUsernameExist struct{}

func (errUsernameExist) Error() string {
	return "Username existed"
}

func (errUsernameExist) StatusCode() int {
	return http.StatusBadRequest
}

type errPhoneNumberExist struct{}

func (errPhoneNumberExist) Error() string {
	return "Phonenumber existed"
}

func (errPhoneNumberExist) StatusCode() int {
	return http.StatusBadRequest
}

type errEmailExist struct{}

func (errEmailExist) Error() string {
	return "Email existed"
}

func (errEmailExist) StatusCode() int {
	return http.StatusBadRequest
}
