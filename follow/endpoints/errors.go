package endpoints

import "net/http"

var (
	//ErrGetFollower ...
	ErrGetFollower = errGetFollower{}
	//ErrGetFollowee ...
	ErrGetFollowee = errGetFollowee{}
)

type errGetFollower struct{}

func (errGetFollower) Error() string {
	return "Get follower error"
}

func (errGetFollower) StatusCode() int {
	return http.StatusBadRequest
}

type errGetFollowee struct{}

func (errGetFollowee) Error() string {
	return "Get followee error"
}

func (errGetFollowee) StatusCode() int {
	return http.StatusBadRequest
}
