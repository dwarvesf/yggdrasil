package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/organization/service"
)

// CreateOrganizationRequest ...
type CreateOrganizationRequest struct {
	Name string `json:"username,omitempty"`
}

// CreateOrganizationResponse ...
type CreateOrganizationResponse struct {
	Token string `json:"token"`
}

// CreateOrganizationEndpoint ...
func CreateOrganizationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return CreateOrganizationResponse{Token: "verifyToken"}, nil
	}
}
