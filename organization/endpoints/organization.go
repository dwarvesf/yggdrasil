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
		orgID, ok := ctx.Value(util.OrgIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

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
}

// ArchiveOrganizationRequest ...
type ArchiveOrganizationRequest struct{}

// ArchiveOrganizationResponse ...
type ArchiveOrganizationResponse struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Status   uint8          `json:"status"`
	Metadata model.Metadata `json:"metadata"`
}

// ArchiveOrganizationEndpoint ...
func ArchiveOrganizationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		orgID, ok := ctx.Value(util.OrgIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

		org := &model.Organization{
			ID:     orgID,
			Status: model.OrganizationStatusInactive,
		}

		org, err := s.OrganizationService.Archive(org)
		if err != nil {
			return nil, err
		}

		return ArchiveOrganizationResponse{
			ID:       org.ID,
			Name:     org.Name,
			Status:   uint8(org.Status),
			Metadata: org.Metadata,
		}, nil
	}
}

// JoinOrganizationRequest ...
type JoinOrganizationRequest struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
}

// JoinOrganizationResponse ...
type JoinOrganizationResponse struct {
	Status string `json:"status"`
}

// JoinOrganizationEndpoint ...
func JoinOrganizationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		orgID, ok := ctx.Value(util.OrgIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

		req := request.(JoinOrganizationRequest)
		uo := &model.UserOrganizations{
			OrganizationID: orgID,
			UserID:         req.UserID,
		}

		err := s.OrganizationService.Join(uo)
		if err != nil {
			return nil, err
		}

		return JoinOrganizationResponse{Status: "success"}, nil
	}
}

// LeaveOrganizationRequest ...
type LeaveOrganizationRequest struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
}

// LeaveOrganizationResponse ...
type LeaveOrganizationResponse struct {
	Status string `json:"status"`
}

// LeaveOrganizationEndpoint ...
func LeaveOrganizationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		orgID, ok := ctx.Value(util.OrgIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

		req := request.(LeaveOrganizationRequest)
		now := time.Now()
		uo := &model.UserOrganizations{
			OrganizationID: orgID,
			UserID:         req.UserID,
			LeftAt:         &now,
		}

		err := s.OrganizationService.Leave(uo)
		if err != nil {
			return nil, err
		}

		return LeaveOrganizationResponse{Status: "success"}, nil
	}
}

// InviteUserOrgRequest ...
type InviteUserOrgRequest struct {
	UserID   uuid.UUID `json:"user_id,omitempty"`
	FromUser uuid.UUID `json:"from_id,omitempty"`
}

// InviteUserOrgResponse ...
type InviteUserOrgResponse struct {
	Status string `json:"status"`
}

// InviteUserOrgEndpoint ...
func InviteUserOrgEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		orgID, ok := ctx.Value(util.OrgIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

		req := request.(InviteUserOrgRequest)
		now := time.Now()
		uo := &model.UserOrganizations{
			OrganizationID: orgID,
			UserID:         req.UserID,
			InvitedBy:      &req.FromUser,
			InvitedAt:      &now,
		}

		err := s.OrganizationService.Join(uo)
		if err != nil {
			return nil, err
		}

		return InviteUserOrgResponse{Status: "success"}, nil
	}
}

// KickUserOrgRequest ...
type KickUserOrgRequest struct {
	UserID   uuid.UUID `json:"user_id,omitempty"`
	FromUser uuid.UUID `json:"from_id,omitempty"`
}

// KickUserOrgResponse ...
type KickUserOrgResponse struct {
	Status string `json:"status"`
}

// KickUserOrgEndpoint ...
func KickUserOrgEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		orgID, ok := ctx.Value(util.OrgIDContextKey).(uuid.UUID)
		if !ok {
			return nil, ErrorMissingID
		}

		req := request.(KickUserOrgRequest)
		now := time.Now()
		uo := &model.UserOrganizations{
			OrganizationID: orgID,
			UserID:         req.UserID,
			KickedBy:       &req.FromUser,
			KickedAt:       &now,
		}

		err := s.OrganizationService.Leave(uo)
		if err != nil {
			return nil, err
		}

		return KickUserOrgResponse{Status: "success"}, nil
	}
}
