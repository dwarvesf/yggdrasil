package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/organization/endpoints"
	"github.com/dwarvesf/yggdrasil/organization/util"
)

// ConfigOrganizationRouter ...
func ConfigOrganizationRouter(r chi.Router, endpoints endpoints.Endpoints, options []httptransport.ServerOption) chi.Router {
	return r.Route("/organizations", func(r chi.Router) {
		r.Post("/", httptransport.NewServer(
			endpoints.CreateOrganization,
			DecodeCreateOrganizationRequest,
			encodeResponse,
			options...,
		).ServeHTTP)
		r.Route("/{orgID}", func(r chi.Router) {
			r.Use(makeOrganizationContext)
			r.Put("/", httptransport.NewServer(
				endpoints.UpdateOrganization,
				DecodeUpdateOrganizationRequest,
				encodeResponse,
				options...,
			).ServeHTTP)
			r.Post("/archive", httptransport.NewServer(
				endpoints.ArchiveOrganization,
				DecodeArchiveOrganizationRequest,
				encodeResponse,
				options...,
			).ServeHTTP)
		})
	})
}

func makeOrganizationContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgID, err := uuid.FromString(chi.URLParam(r, "orgID"))

		if err != nil {
			http.Error(w, `{"error": "INVALID_ID"}`, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), util.OrgIDContextKey, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
