package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/dwarvesf/yggdrasil/services/device/endpoints"
	"github.com/dwarvesf/yggdrasil/services/device/service"
)

//NewHTTPHandler create new http handler
func NewHTTPHandler(s service.Service, endpoints endpoints.Endpoints, logger log.Logger, useCORS bool) http.Handler {
	r := chi.NewRouter()
	if useCORS {
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Acceot", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
		})
		r.Use(cors.Handler)
	}

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Route("/", func(r chi.Router) {
		r.Post("/", httptransport.NewServer(
			endpoints.CreateDevice,
			DecodeCreateDeviceRequest,
			encodeResponse,
			options...,
		).ServeHTTP)
		r.Get("/{device_id}", httptransport.NewServer(
			endpoints.GetDevice,
			DecodeGetDeviceRequest,
			encodeResponse,
			options...,
		).ServeHTTP)
		r.Get("/user/{user_id}", httptransport.NewServer(
			endpoints.GetListDevice,
			DecodeGetListDeviceRequest,
			encodeResponse,
			options...,
		).ServeHTTP)
		r.Put("/{device_id}", httptransport.NewServer(
			endpoints.UpdateDevice,
			DecodeUpdateDeviceRequest,
			encodeResponse,
			options...,
		).ServeHTTP)
	})

	return r
}
