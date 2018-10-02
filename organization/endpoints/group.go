package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/organization/model"
	"github.com/dwarvesf/yggdrasil/organization/service"
	"github.com/dwarvesf/yggdrasil/organization/util"
)

// CreateGroupRequest ...
type CreateGroupRequest struct {
	Name           string         `json:"name,omitempty"`
	OrganizationID uuid.UUID      `json:"organization_id,omitempty"`
	Metadata       model.Metadata `json:"metadata,omitempty"`
}

// CreateGroupResponse ...
type CreateGroupResponse struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Status   uint8          `json:"status"`
	Metadata model.Metadata `json:"metadata"`
}

// CreateGroupEndpoint ...
func CreateGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateGroupRequest)
		g := &model.Group{
			Name:           req.Name,
			Metadata:       req.Metadata,
			OrganizationID: req.OrganizationID,
		}

		g, err := s.GroupService.Create(g)
		if err != nil {
			return nil, err
		}

		return CreateGroupResponse{
			ID:       g.ID,
			Name:     g.Name,
			Status:   uint8(g.Status),
			Metadata: g.Metadata,
		}, nil
	}
}

// UpdateGroupRequest ...
type UpdateGroupRequest struct {
	Name     string         `json:"name,omitempty"`
	Metadata model.Metadata `json:"metadata,omitempty"`
}

// UpdateGroupResponse ...
type UpdateGroupResponse struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Status   uint8          `json:"status"`
	Metadata model.Metadata `json:"metadata"`
}

// UpdateGroupEndpoint ...
func UpdateGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID); ok {
			req := request.(UpdateGroupRequest)
			g := &model.Group{
				ID:       groupID,
				Name:     req.Name,
				Metadata: req.Metadata,
			}

			g, err := s.GroupService.Update(g)
			if err != nil {
				return nil, err
			}

			return UpdateGroupResponse{
				ID:       g.ID,
				Name:     g.Name,
				Status:   uint8(g.Status),
				Metadata: g.Metadata,
			}, nil
		}

		return nil, ErrorMissingID
	}
}
