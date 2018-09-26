package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dwarvesf/yggdrasil/follow/endpoints"
)

// DecodeCreateOrganizationRequest ...
func DecodeCreateFollowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateFollowRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
