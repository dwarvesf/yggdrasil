package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/networks/model"
	"github.com/dwarvesf/yggdrasil/services/networks/service"
)

//MakeFriendRequest use for api Makefriend
type MakeFriendRequest struct {
	FromUser uuid.UUID `json:"from"`
	ToUser   uuid.UUID `json:"to"`
}

//MakeFriendResponse use for api Makefriend
type MakeFriendResponse struct {
	Status string `json:"status"`
}

//MakeFriendEndpoints endpoint for makefriend request
func MakeFriendEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MakeFriendRequest)

		err := s.FriendService.MakeFriend(req.FromUser, req.ToUser)
		if err != nil {
			return nil, err
		}

		return MakeFriendResponse{Status: "success"}, nil
	}
}

//AcceptRequest use for api Accept request
type AcceptRequest struct {
	FromUser uuid.UUID `json:"from"`
	ToUser   uuid.UUID `json:"to"`
}

//AcceptResponse response api Accept request
type AcceptResponse struct {
	Status string `json:"status"`
}

//AcceptRequestEndpoints endpoint for accept request
func AcceptRequestEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AcceptRequest)

		err := s.FriendService.Accept(req.FromUser, req.ToUser)
		if err != nil {
			return nil, err
		}

		return AcceptResponse{Status: "accepted"}, nil
	}
}

//RejectRequest use for reject requenst
type RejectRequest struct {
	FromUser uuid.UUID `json:"from"`
	ToUser   uuid.UUID `json:"to"`
}

//RejectResponse respone for reject request
type RejectResponse struct {
	Status string `json:"status"`
}

//RejectRequestEndpoints endpoint for reject request
func RejectRequestEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RejectRequest)

		err := s.FriendService.Reject(req.FromUser, req.ToUser)
		if err != nil {
			return nil, err
		}

		return RejectResponse{Status: "rejected"}, nil
	}
}

//UnFriendtRequest use for unfriend request
type UnFriendtRequest struct {
	FromUser uuid.UUID `json:"from"`
	ToUser   uuid.UUID `json:"to"`
}

//UnFriendResponse respone for unfreind request
type UnFriendResponse struct {
	Status string `json:"status"`
}

//UnFriendEndpoints endpoint for Unfreind request
func UnFriendEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UnFriendtRequest)

		err := s.FriendService.UnFriend(req.FromUser, req.ToUser)
		if err != nil {
			return nil, err
		}

		return UnFriendResponse{Status: "success"}, nil
	}
}

//GetFriendRequest use for get friend request, with friend's id is userID
type GetFriendRequest struct {
	UserID uuid.UUID `json:"userID"`
}

//GetFriendResponse response for getfriend response
type GetFriendResponse struct {
	Friends []model.Friend `json:"friends"`
}

//GetFriendsEndpoints endpoint for get friend request
func GetFriendsEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetFriendRequest)

		res, err := s.FriendService.GetFriends(req.UserID)
		if err != nil {
			return nil, err
		}

		return GetFriendResponse{Friends: res}, nil
	}
}
