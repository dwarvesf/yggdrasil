package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/organization/service"
)

type Endpoints struct {
	CreateOrganization endpoint.Endpoint
	UpdateOrganization endpoint.Endpoint
	CreateGroup        endpoint.Endpoint
	UpdateGroup        endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreateOrganization: CreateOrganizationEndpoint(s),
		UpdateOrganization: UpdateOrganizationEndpoint(s),
		CreateGroup:        CreateGroupEndpoint(s),
		UpdateGroup:        UpdateGroupEndpoint(s),
	}
}
