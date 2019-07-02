package follow

import "net/http"

var (
	//ErrFromUserUUIDEmpty ...
	ErrFromUserUUIDEmpty = errFromUserUUIDEmpty{}
	//ErrToUserUUIDEmpty ...
	ErrToUserUUIDEmpty = errToUserUUIDEmpty{}
	//ErrorCreateFollow ...
	ErrorCreateFollow = errCreateFollow{}
	//ErrUnfollow ...
	ErrUnfollow = errUnfollow{}
)

type errFromUserUUIDEmpty struct{}

func (errFromUserUUIDEmpty) Error() string {
	return "FromUser UUID empty"
}

func (errFromUserUUIDEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errToUserUUIDEmpty struct{}

func (errToUserUUIDEmpty) Error() string {
	return "ToUser UUID empty"
}

func (errToUserUUIDEmpty) StatusCode() int {
	return http.StatusBadRequest
}

type errCreateFollow struct{}

func (errCreateFollow) Error() string {
	return "Follow error"
}

func (errCreateFollow) StatusCode() int {
	return http.StatusBadRequest
}

type errUnfollow struct{}

func (errUnfollow) Error() string {
	return "User hasn't existed yet"
}

func (errUnfollow) StatusCode() int {
	return http.StatusBadRequest
}
