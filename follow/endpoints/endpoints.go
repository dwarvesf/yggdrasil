package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/follow/service"
)

type Endpoints struct {
	CreateFollow endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreateFollow: CreateFollowEndpoint(s),
	}
}
