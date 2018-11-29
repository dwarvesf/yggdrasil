package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/dwarvesf/yggdrasil/services/identity/endpoints"
	"github.com/dwarvesf/yggdrasil/services/identity/service"
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
	r.Route("/users", func(r chi.Router) {
		r.Get("/{id}", httptransport.NewServer(
			endpoints.GetUser,
			DecodeGetUserRequest,
			encodeResponse,
			options...,
		).ServeHTTP)
	})

	r.Post("/register", httptransport.NewServer(
		endpoints.Register,
		DecodeCreateUserRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/user_verify", httptransport.NewServer(
		endpoints.VerifyUser,
		DecodeVerifyUserRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/authenticate", httptransport.NewServer(
		endpoints.VerifyToken,
		DecodeVerifyTokenRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/login", httptransport.NewServer(
		endpoints.Login,
		DecodeLoginRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Post("/referral", httptransport.NewServer(
		endpoints.ReferralUser,
		DecodeReferralRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Put("/referral/response", httptransport.NewServer(
		endpoints.ReferralResponse,
		DecodeReferralResponse,
		encodeResponse,
		options...,
	).ServeHTTP)

	return r
}
