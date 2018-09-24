package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/dwarvesf/yggdrasil/organization/endpoints"
	"github.com/dwarvesf/yggdrasil/organization/service"
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

	//TODO: Add authorization for all endpoints

	// endpoints
	r.Post("/", httptransport.NewServer(
		endpoints.CreateOrganization,
		DecodeCreateOrganizationRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	return r
}
