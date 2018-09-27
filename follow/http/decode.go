package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/follow/endpoints"
)

//DecodeCreateFollowRequest ...
func DecodeCreateFollowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateFollowRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

//DecodeUnFollowRequest ...
func DecodeUnFollowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UnFollowRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

//DecodeFollowerListRequest ...
func DecodeFollowerListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := uuid.FromString(chi.URLParam(r, "user_id"))
	if err != nil {
		return nil, err
	}

	return endpoints.FollowerListRequest{UserID: userID}, nil
}

//DecodeFolloweeListRequest ...
func DecodeFolloweeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID, err := uuid.FromString(chi.URLParam(r, "user_id"))
	if err != nil {
		return nil, err
	}

	return endpoints.FolloweeListRequest{UserID: userID}, nil
}
