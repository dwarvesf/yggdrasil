package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"

	cfg "github.com/dwarvesf/yggdrasil/organization/cmd/config"
	"github.com/dwarvesf/yggdrasil/organization/db"
	"github.com/dwarvesf/yggdrasil/organization/endpoints"
	serviceHttp "github.com/dwarvesf/yggdrasil/organization/http"
	"github.com/dwarvesf/yggdrasil/organization/middlewares"
	"github.com/dwarvesf/yggdrasil/organization/service"
	"github.com/dwarvesf/yggdrasil/organization/service/organization"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

func main() {
	var (
		httpAddr = flag.String("addr", fmt.Sprintf(":%v", os.Getenv("PORT")), "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("consul-server:8500"),
	})
	if err != nil {
		panic(err)
	}

	// Get jwt secret
	secretKey, secretKeyErr := toolkit.GetConsulValueFromKey(consulClient, "jwt_secret")
	if secretKeyErr != nil {
		panic(secretKeyErr)
	}
	cfg.JwtSecret = secretKey

	pgdb, closeDB := db.New(consulClient)
	db.Migrate(pgdb)
	defer closeDB()

	var s service.Service
	{
		s = service.Service{
			OrganizationService: middlewares.Compose(
				organization.NewPGService(pgdb),
				organization.ValidationMiddleware(),
			).(organization.Service),
		}
	}

	var h http.Handler
	{
		h = serviceHttp.NewHTTPHandler(
			s,
			endpoints.MakeServerEndpoints(s),
			logger,
			true,
		)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}

		agent := consulClient.Agent()

		name := "organization"
		if err := agent.ServiceRegister(&consul.AgentServiceRegistration{
			Name:    name,
			Port:    port,
			Address: os.Getenv("PRIVATE_IP"),
		}); err != nil {
			panic(err)
		}
		logger.Log("consul", "registered", "name", name)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
