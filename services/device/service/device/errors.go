package device

import "net/http"

//Error declaration
var (
	ErrInvalidUser        = errInvalidUser{}
	ErrInvalidType        = errInvalidType{}
	ErrUserIDIsRequired   = errUserIDIsRequired{}
	ErrDeviceIDIsRequired = errDeviceIDIsRequired{}
	ErrDeviceNotExist     = errDeviceNotExist{}
	ErrInvalidStatus      = errInvalidStatus{}
)

type errInvalidUser struct{}

func (errInvalidUser) Error() string {
	return "Invalid user"
}

func (errInvalidUser) StatusCode() int {
	return http.StatusNonAuthoritativeInfo
}

type errInvalidType struct{}

func (errInvalidType) Error() string {
	return "Invalid device type"
}

func (errInvalidType) StatusCode() int {
	return http.StatusBadRequest
}

type errUserIDIsRequired struct{}

func (errUserIDIsRequired) Error() string {
	return "User ID is required"
}

func (errUserIDIsRequired) StatusCode() int {
	return http.StatusBadRequest
}

type errDeviceIDIsRequired struct{}

func (errDeviceIDIsRequired) Error() string {
	return "Device ID is required"
}

func (errDeviceIDIsRequired) StatusCode() int {
	return http.StatusBadRequest
}

type errDeviceNotExist struct{}

func (errDeviceNotExist) Error() string {
	return "Device not exist"
}

func (errDeviceNotExist) StatusCode() int {
	return http.StatusBadRequest
}

type errInvalidStatus struct{}

func (errInvalidStatus) Error() string {
	return "Invalid device status"
}

func (errInvalidStatus) StatusCode() int {
	return http.StatusBadRequest
}
