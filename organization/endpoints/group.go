package endpoints

import (
	"context"
	"time"

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

// JoinGroupRequest ...
type JoinGroupRequest struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
}

// JoinGroupResponse ...
type JoinGroupResponse struct {
	Status string `json:"status"`
}

// JoinGroupEndpoint ...
func JoinGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID); ok {
			req := request.(JoinGroupRequest)
			ug := &model.UserGroups{
				GroupID: groupID,
				UserID:  req.UserID,
			}

			err := s.GroupService.Join(ug)
			if err != nil {
				return nil, err
			}

			return JoinGroupResponse{Status: "success"}, nil
		}

		return nil, ErrorMissingID
	}
}

// LeaveGroupRequest ...
type LeaveGroupRequest struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
}

// LeaveGroupResponse ...
type LeaveGroupResponse struct {
	Status string `json:"status"`
}

// LeaveGroupEndpoint ...
func LeaveGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID); ok {
			req := request.(LeaveGroupRequest)
			now := time.Now()
			ug := &model.UserGroups{
				GroupID: groupID,
				UserID:  req.UserID,
				LeftAt:  &now,
			}

			err := s.GroupService.Leave(ug)
			if err != nil {
				return nil, err
			}

			return LeaveGroupResponse{Status: "success"}, nil
		}

		return nil, ErrorMissingID
	}
}

// InviteUserRequest ...
type InviteUserRequest struct {
	UserID   uuid.UUID `json:"user_id,omitempty"`
	FromUser uuid.UUID `json:"from_id,omitempty"`
}

// InviteUserResponse ...
type InviteUserResponse struct {
	Status string `json:"status"`
}

// InviteUserEndpoint ...
func InviteUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID); ok {
			req := request.(InviteUserRequest)
			now := time.Now()
			ug := &model.UserGroups{
				GroupID:   groupID,
				UserID:    req.UserID,
				InvitedBy: &req.FromUser,
				InvitedAt: &now,
			}

			err := s.GroupService.Join(ug)
			if err != nil {
				return nil, err
			}

			return InviteUserResponse{Status: "success"}, nil
		}

		return nil, ErrorMissingID
	}
}

// KickUserRequest ...
type KickUserRequest struct {
	UserID   uuid.UUID `json:"user_id,omitempty"`
	FromUser uuid.UUID `json:"from_id,omitempty"`
}

// KickUserResponse ...
type KickUserResponse struct {
	Status string `json:"status"`
}

// KickUserEndpoint ...
func KickUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID); ok {
			req := request.(KickUserRequest)
			now := time.Now()
			ug := &model.UserGroups{
				GroupID:  groupID,
				UserID:   req.UserID,
				KickedBy: &req.FromUser,
				KickedAt: &now,
			}

			err := s.GroupService.Leave(ug)
			if err != nil {
				return nil, err
			}

			return KickUserResponse{Status: "success"}, nil
		}

		return nil, ErrorMissingID
	}
}
