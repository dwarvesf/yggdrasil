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
