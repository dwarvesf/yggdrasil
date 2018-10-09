package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dwarvesf/yggdrasil/organization/endpoints"
)

// DecodeCreateOrganizationRequest ...
func DecodeCreateOrganizationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateOrganizationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeUpdateOrganizationRequest ...
func DecodeUpdateOrganizationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateOrganizationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeArchiveOrganizationRequest ...
func DecodeArchiveOrganizationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ArchiveGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeCreateGroupRequest ...
func DecodeCreateGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeUpdateGroupRequest ...
func DecodeUpdateGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeArchiveGroupRequest ...
func DecodeArchiveGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ArchiveGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeJoinGroupRequest ...
func DecodeJoinGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.JoinGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeLeaveGroupRequest ...
func DecodeLeaveGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LeaveGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeInviteUserRequest ...
func DecodeInviteUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.InviteUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeKickUserRequest ...
func DecodeKickUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.KickUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeJoinOrgRequest ...
func DecodeJoinOrgRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.JoinOrganizationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeLeaveOrgRequest ...
func DecodeLeaveOrgRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LeaveOrganizationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeInviteUserOrgRequest ...
func DecodeInviteUserOrgRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.InviteUserOrgRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// DecodeKickUserOrgRequest ...
func DecodeKickUserOrgRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.KickUserOrgRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
