package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/organization/service"
)

type Endpoints struct {
	CreateOrganization     endpoint.Endpoint
	UpdateOrganization     endpoint.Endpoint
	ArchiveOrganization    endpoint.Endpoint
	JoinOrganization       endpoint.Endpoint
	LeaveOrganization      endpoint.Endpoint
	InviteUserOrganization endpoint.Endpoint
	KickUserOrganization   endpoint.Endpoint
	CreateGroup            endpoint.Endpoint
	UpdateGroup            endpoint.Endpoint
	ArchiveGroup           endpoint.Endpoint
	JoinGroup              endpoint.Endpoint
	LeaveGroup             endpoint.Endpoint
	InviteUser             endpoint.Endpoint
	KickUser               endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		CreateOrganization:     CreateOrganizationEndpoint(s),
		UpdateOrganization:     UpdateOrganizationEndpoint(s),
		ArchiveOrganization:    ArchiveOrganizationEndpoint(s),
		JoinOrganization:       JoinOrganizationEndpoint(s),
		LeaveOrganization:      LeaveOrganizationEndpoint(s),
		InviteUserOrganization: InviteUserOrgEndpoint(s),
		KickUserOrganization:   KickUserOrgEndpoint(s),
		CreateGroup:            CreateGroupEndpoint(s),
		UpdateGroup:            UpdateGroupEndpoint(s),
		ArchiveGroup:           ArchiveGroupEndpoint(s),
		JoinGroup:              JoinGroupEndpoint(s),
		LeaveGroup:             LeaveGroupEndpoint(s),
		InviteUser:             InviteUserEndpoint(s),
		KickUser:               KickUserEndpoint(s),
	}
}
