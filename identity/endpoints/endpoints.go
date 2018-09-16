package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/identity/service"
)

type Endpoints struct {
	GetUser endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		GetUser: MakeGetUserEndpoint(s),
	}
}
