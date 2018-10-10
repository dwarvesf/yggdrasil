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
		groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

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
}

// ArchiveGroupResponse ...
type ArchiveGroupResponse struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Status   uint8          `json:"status"`
	Metadata model.Metadata `json:"metadata"`
}

// ArchiveGroupEndpoint ...
func ArchiveGroupEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

		g := &model.Group{
			ID:     groupID,
			Status: model.GroupStatusInactive,
		}

		g, err := s.GroupService.Archive(g)
		if err != nil {
			return nil, err
		}

		return ArchiveGroupResponse{
			ID:       g.ID,
			Name:     g.Name,
			Status:   uint8(g.Status),
			Metadata: g.Metadata,
		}, nil
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
		groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

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
		groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

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
		groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

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
		groupID, ok := ctx.Value(util.GroupIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

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
}
