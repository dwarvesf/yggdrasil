package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dwarvesf/yggdrasil/identity/endpoints"
	"github.com/go-chi/chi"
)

// DecodeGetUserRequest ...
func DecodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.GetUserRequest{ID: chi.URLParam(r, "id")}, nil
}

//DecodeCreateUserRequest ...
func DecodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

//DecodeVerifyUserRequest decode verify request
func DecodeVerifyUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.VerifyUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

//DecodeLoginRequest ...
func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

//DecodeVerifyTokenRequest ...
func DecodeVerifyTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.VerifyTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
