package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/dwarvesf/yggdrasil/networks/endpoints"
	"github.com/dwarvesf/yggdrasil/networks/service"
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

	r.Route("/", func(r chi.Router) {
		r.Post("/follow", httptransport.NewServer(
			endpoints.CreateFollow,
			DecodeCreateFollowRequest,
			encodeResponse,
			options...,
		).ServeHTTP)

		r.Put("/unfollow", httptransport.NewServer(
			endpoints.UnFollow,
			DecodeUnFollowRequest,
			encodeResponse,
			options...,
		).ServeHTTP)
	})

	r.Get("/follower/{user_id}", httptransport.NewServer(
		endpoints.GetFollower,
		DecodeFollowerListRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Get("/followee/{user_id}", httptransport.NewServer(
		endpoints.GetFollowee,
		DecodeFolloweeListRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	return r
}
