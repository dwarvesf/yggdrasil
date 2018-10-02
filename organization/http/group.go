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

// ConfigGroupRouter used to declare all group endpoints
func ConfigGroupRouter(r chi.Router, endpoints endpoints.Endpoints, options []httptransport.ServerOption) chi.Router {
	return r.Route("/groups", func(r chi.Router) {
		r.Post("/", httptransport.NewServer(
			endpoints.CreateGroup,
			DecodeCreateGroupRequest,
			encodeResponse,
			options...,
		).ServeHTTP)
		r.Route("/{groupID}", func(r chi.Router) {
			r.Use(makeGroupContext)
			r.Put("/", httptransport.NewServer(
				endpoints.UpdateGroup,
				DecodeUpdateGroupRequest,
				encodeResponse,
				options...,
			).ServeHTTP)
		})
	})
}

func makeGroupContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupID, err := uuid.FromString(chi.URLParam(r, "groupID"))

		if err != nil {
			http.Error(w, `{"error": "INVALID_ID"}`, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), util.GroupIDContextKey, groupID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
