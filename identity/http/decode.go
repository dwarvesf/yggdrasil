package http

import (
	"context"
	"net/http"

	"github.com/dwarvesf/yggdrasil/identity/endpoints"
	"github.com/go-chi/chi"
)

// DecodeGetUserRequest ...
func DecodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.GetUserRequest{ID: chi.URLParam(r, "id")}, nil
}
