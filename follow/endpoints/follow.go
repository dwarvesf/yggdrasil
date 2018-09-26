package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/follow/service"
)

// CreateFollowRequest ...
type CreateFollowRequest struct {
	Name string `json:"username,omitempty"`
}

// CreateFollowResponse ...
type CreateFollowResponse struct {
	Token string `json:"token"`
}

// CreateFollowEndpoint ...
func CreateFollowEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return CreateFollowResponse{Token: "verifyToken"}, nil
	}
}
