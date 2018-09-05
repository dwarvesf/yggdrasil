package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/dwarvesf/yggdrasil/identity/endpoints"
	"github.com/dwarvesf/yggdrasil/identity/service"
)

// NewHTTPHandler ...
func NewHTTPHandler(s service.Service, endpoints endpoints.Endpoints, logger log.Logger, useCORS bool) http.Handler {
	r := chi.NewRouter()

	if useCORS {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
		})
		r.Use(cors.Handler)
	}

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// endpoints
	r.Post("/add", httptransport.NewServer(
		endpoints.Add,
		DecodeAddRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	return r
}

// DecodeAddRequest ...
func DecodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
