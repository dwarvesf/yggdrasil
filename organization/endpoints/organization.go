package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/organization/service"
)

// CreateOrganizationRequest ...
type CreateOrganizationRequest struct {
	Name string `json:"name,omitempty"`
}

// CreateOrganizationResponse ...
type CreateOrganizationResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status uint8     `json:"status"`
}

// CreateOrganizationEndpoint ...
func CreateOrganizationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateOrganizationRequest)

		org, err := s.OrganizationService.Create(req.Name)
		if err != nil {
			return nil, err
		}

		return CreateOrganizationResponse{
			ID:     org.ID,
			Name:   org.Name,
			Status: uint8(org.Status),
		}, nil
	}
}

// UpdateOrganizationRequest ...
type UpdateOrganizationRequest struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name,omitempty"`
}

// UpdateOrganizationResponse ...
type UpdateOrganizationResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status uint8     `json:"status"`
}

// UpdateOrganizationEndpoint ...
func UpdateOrganizationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateOrganizationRequest)

		org, err := s.OrganizationService.Update(req.ID, req.Name)
		if err != nil {
			return nil, err
		}

		return UpdateOrganizationResponse{
			ID:     org.ID,
			Name:   org.Name,
			Status: uint8(org.Status),
		}, nil
	}
}
