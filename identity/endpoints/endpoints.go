package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/identity/service"
)

type Endpoints struct {
	GetUser          endpoint.Endpoint
	Register         endpoint.Endpoint
	VerifyUser       endpoint.Endpoint
	Login            endpoint.Endpoint
	VerifyToken      endpoint.Endpoint
	ReferralUser     endpoint.Endpoint
	ReferralResponse endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		GetUser:          MakeGetUserEndpoint(s),
		Register:         RegisterEndpoint(s),
		VerifyUser:       VerifyUserEndpoints(s),
		Login:            LoginEndpoint(s),
		VerifyToken:      VerifyTokenEndpoints(s),
		ReferralUser:     ReferralUser(s),
		ReferralResponse: ResponseReferral(s),
	}
}
