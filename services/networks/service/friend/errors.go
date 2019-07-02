package friend

import (
	"net/http"
)

var (
	//ErrFromUserEmpty ...
	ErrFromUserEmpty = errFromUserEmpty{}
	//ErrToUserEmpty ...
	ErrToUserEmpty = errToUserEmpty{}
	//ErrRequestNotExist ...
	ErrRequestNotExist = errRequestNotExist{}
	//ErrGetCondition ...
	ErrGetCondition = errGetCondition{}
	//ErrUserIDInvalid ...
	ErrUserIDInvalid = errUserIDInvalid{}
	//ErrFriendNotAccepted ...
	ErrFriendAccepted = errFriendAccepted{}
	//ErrFriendRejected ....
	ErrFriendRejected = errFriendRejected{}
)

type errFromUserEmpty struct{}

func (errFromUserEmpty) Error() string {
	return "From UserID empty"
}

func (errFromUserEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errToUserEmpty struct{}

func (errToUserEmpty) Error() string {
	return "To UserID empty"
}

func (errToUserEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errRequestNotExist struct{}

func (errRequestNotExist) Error() string {
	return "Relation ship has not existed yet"
}

func (errRequestNotExist) StatusCode() int {
	return http.StatusBadRequest
}

type errGetCondition struct{}

func (errGetCondition) Error() string {
	return "Get condition error"
}

func (errGetCondition) StatusCode() int {
	return http.StatusBadRequest
}

type errUserIDInvalid struct{}

func (errUserIDInvalid) Error() string {
	return "UserID invalid"
}

func (errUserIDInvalid) StatusCode() int {
	return http.StatusBadRequest
}

type errFriendAccepted struct{}

func (errFriendAccepted) Error() string {
	return "Friend accepted"
}

func (errFriendAccepted) StatusCode() int {
	return http.StatusBadRequest
}

type errFriendRejected struct{}

func (errFriendRejected) Error() string {
	return "Friend not accepted"
}

func (errFriendRejected) StatusCode() int {
	return http.StatusBadRequest
}
