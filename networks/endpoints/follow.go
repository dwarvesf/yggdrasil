package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/networks/model"
	"github.com/dwarvesf/yggdrasil/networks/service"
	"github.com/dwarvesf/yggdrasil/networks/service/follow"
)

// CreateFollowRequest ...
type CreateFollowRequest struct {
	FromUser uuid.UUID `json:"from"`	
	ToUser   uuid.UUID `json:"to"`
}

//CreateFollowResponse ...
type CreateFollowResponse struct {
	Status string `json:"status"`
}

//CreateFollowEndpoints ...
func CreateFollowEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateFollowRequest)

		err := s.FollowService.Follow(req.FromUser, req.ToUser)
		if err != nil {
			return nil, err
		}
		return CreateFollowResponse{Status: "Following"}, nil
	}
}

//UnFollowRequest ...
type UnFollowRequest struct {
	FromUser uuid.UUID `json:"from"`
	ToUser   uuid.UUID `json:"to"`
}

//UnFollowResponse ...
type UnFollowResponse struct {
	Status string `json:"status"`
}

//UnFollowEndpoints ...
func UnFollowEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UnFollowRequest)

		err := s.FollowService.UnFollow(req.FromUser, req.ToUser)
		if err != nil {
			return nil, err
		}
		return UnFollowResponse{Status: "UnFollowed"}, nil
	}
}

//FollowerListRequest ...
type FollowerListRequest struct {
	UserID uuid.UUID
}

//FollowerListResponse ...
type FollowerListResponse struct {
	User []model.Follow `json:"user"`
}

//MakeFollowerEndpoints get all user, you are following
func MakeFollowerEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FollowerListRequest)
		res, err := s.FollowService.FindAll(&follow.Query{ToUser: req.UserID, Status: 1})
		if err != nil {
			return nil, ErrGetFollower
		}
		return res, nil
	}
}

//FolloweeListRequest ...
type FolloweeListRequest struct {
	UserID uuid.UUID
}

//FolloweeListResponse ...
type FolloweeListResponse struct {
	User []model.Follow `json:"user"`
}

//MakeFolloweeEndpoints get all user are following you
func MakeFolloweeEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FolloweeListRequest)
		res, err := s.FollowService.FindAll(&follow.Query{FromUser: req.UserID, Status: 1})
		if err != nil {
			return nil, ErrGetFollowee
		}
		return res, nil
	}
}
