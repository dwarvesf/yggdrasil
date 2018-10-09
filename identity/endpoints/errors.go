package endpoints

import "net/http"

var (
	//ErrorInvalidLogin return login fail message
	ErrorInvalidLogin = errInvalidLogin{}
	//ErrorNotActive return token invalid message
	ErrorNotActive = errNotActive{}
	//ErrUnauthorize ...
	ErrUnauthorize = errUnauthorize{}
	//ErrReferralUserInfo userinfo error
	ErrReferralUserInfo = errReferralUserInfo{}
	//ErrCreateNewAccount create new account error
	ErrCreateNewAccount = errCreateNewAccount{}
	//ErrTTLExpires ....
	ErrTTLExpires = errTTLExpires{}
	//ErrDeleteRefferal ...
	ErrDeleteRefferal = errDeleteRefferal{}
	//ErrCodeInvalid ...
	ErrCodeInvalid = errCodeInvalid{}
)

type errReferralUserInfo struct{}

func (errReferralUserInfo) Error() string {
	return "User info error"
}

func (errReferralUserInfo) StatusCode() int {
	return http.StatusBadRequest
}

type errCreateNewAccount struct{}

func (errCreateNewAccount) Error() string {
	return "Create new account error"
}

func (errCreateNewAccount) StatusCode() int {
	return http.StatusBadRequest
}

type errTTLExpires struct{}

func (errTTLExpires) Error() string {
	return "TTL expires"
}

func (errTTLExpires) StatusCode() int {
	return http.StatusBadRequest
}

type errCodeInvalid struct{}

func (errCodeInvalid) Error() string {
	return "Code in valid"
}

func (errCodeInvalid) StatusCode() int {
	return http.StatusBadRequest
}

type errDeleteRefferal struct{}

func (errDeleteRefferal) Error() string {
	return "Delete referral error"
}

func (errDeleteRefferal) StatusCode() int {
	return http.StatusBadRequest
}

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
