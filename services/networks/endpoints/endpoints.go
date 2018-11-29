package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/services/networks/service"
)

//Endpoints ...
type Endpoints struct {
	CreateFollow endpoint.Endpoint
	UnFollow     endpoint.Endpoint
	GetFollower  endpoint.Endpoint
	GetFollowee  endpoint.Endpoint

	MakeFriend endpoint.Endpoint
	Accept     endpoint.Endpoint
	Reject     endpoint.Endpoint
	UnFriend   endpoint.Endpoint
	GetFriends endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreateFollow: CreateFollowEndpoints(s),
		UnFollow:     UnFollowEndpoints(s),
		GetFollower:  MakeFollowerEndpoints(s),
		GetFollowee:  MakeFolloweeEndpoints(s),

		MakeFriend: MakeFriendEndpoints(s),
		Accept:     AcceptRequestEndpoints(s),
		Reject:     RejectRequestEndpoints(s),
		UnFriend:   UnFriendEndpoints(s),
		GetFriends: GetFriendsEndpoints(s),
	}
}
