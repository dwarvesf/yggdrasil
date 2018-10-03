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
