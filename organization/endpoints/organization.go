package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/organization/model"
	"github.com/dwarvesf/yggdrasil/organization/service"
	"github.com/dwarvesf/yggdrasil/organization/util"
)

// CreateOrganizationRequest ...
type CreateOrganizationRequest struct {
	Name     string         `json:"name,omitempty"`
	Metadata model.Metadata `json:"metadata,omitempty"`
}

// CreateOrganizationResponse ...
type CreateOrganizationResponse struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Status   uint8          `json:"status"`
	Metadata model.Metadata `json:"metadata"`
}

// CreateOrganizationEndpoint ...
func CreateOrganizationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateOrganizationRequest)
		org := &model.Organization{
			Name:     req.Name,
			Metadata: req.Metadata,
		}

		org, err := s.OrganizationService.Create(org)
		if err != nil {
			return nil, err
		}

		return CreateOrganizationResponse{
			ID:       org.ID,
			Name:     org.Name,
			Status:   uint8(org.Status),
			Metadata: org.Metadata,
		}, nil
	}
}

// UpdateOrganizationRequest ...
type UpdateOrganizationRequest struct {
	Name     string         `json:"name,omitempty"`
	Metadata model.Metadata `json:"metadata,omitempty"`
}

// UpdateOrganizationResponse ...
type UpdateOrganizationResponse struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Status   uint8          `json:"status"`
	Metadata model.Metadata `json:"metadata"`
}

// UpdateOrganizationEndpoint ...
func UpdateOrganizationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if orgID, ok := ctx.Value(util.OrgIDContextKey).(uuid.UUID); ok {
			req := request.(UpdateOrganizationRequest)
			org := &model.Organization{
				ID:       orgID,
				Name:     req.Name,
				Metadata: req.Metadata,
			}

			org, err := s.OrganizationService.Update(org)
			if err != nil {
				return nil, err
			}

			return UpdateOrganizationResponse{
				ID:       org.ID,
				Name:     org.Name,
				Status:   uint8(org.Status),
				Metadata: org.Metadata,
			}, nil
		}

		return nil, ErrorMissingID
	}
}
