package endpoints

import (
	"context"

	"github.com/k0kubun/pp"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/identity/service"
)

type GetUserRequest struct {
	ID string `json:"id"`
}

// MakeGetUserEndpoint ...
func MakeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		pp.Println(request.(GetUserRequest))
		return s.UserService.Get(request.(GetUserRequest).ID)
	}
}
