package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/yggdrasil/services/organization/endpoints"
	"github.com/dwarvesf/yggdrasil/services/organization/middlewares"
	"github.com/dwarvesf/yggdrasil/services/organization/service"
	"github.com/dwarvesf/yggdrasil/services/organization/service/group"
	"github.com/dwarvesf/yggdrasil/services/organization/service/organization"
)

// NewHTTPHandler that create main handler of the app
func NewHTTPHandler(pgdb *gorm.DB, logger log.Logger, useCORS bool) http.Handler {
	s := service.Service{
		OrganizationService: middlewares.Compose(
			organization.NewPGService(pgdb),
			organization.ValidationMiddleware(),
		).(organization.Service),
		GroupService: middlewares.Compose(
			group.NewPGService(pgdb),
			group.ValidationMiddleware(),
		).(group.Service),
	}

	return configHandler(s, endpoints.MakeServerEndpoints(s), logger, useCORS)
}

func configHandler(s service.Service, endpoints endpoints.Endpoints, logger log.Logger, useCORS bool) http.Handler {
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

	// TODO: Add authorization for all endpoints

	ConfigOrganizationRouter(r, endpoints, options)
	ConfigGroupRouter(r, endpoints, options)

	return r
}
