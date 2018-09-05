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

	"github.com/dwarvesf/yggdrasil/identity/endpoints"
	serviceHttp "github.com/dwarvesf/yggdrasil/identity/http"
	"github.com/dwarvesf/yggdrasil/identity/middlewares"
	"github.com/dwarvesf/yggdrasil/identity/postgres"
	"github.com/dwarvesf/yggdrasil/identity/service"
	"github.com/dwarvesf/yggdrasil/identity/service/add"
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

	// FIXME: replace this with `postgres.New()`
	pgdb, close := postgres.NewFake(os.Getenv("PG_DATASOURCE"))
	defer func() {
		if err := close(); err != nil {
			logger.Log("msg", "failed to close postgres connection", "err", err)
		}
	}()

	var s service.Service
	{
		s = service.Service{
			AddService: middlewares.Compose(
				postgres.NewAddStore(pgdb),
				add.LoggingMiddleware(logger),
				add.ValidationMiddleware(),
			).(add.Service),
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

		client, err := consul.NewClient(&consul.Config{
			Address: fmt.Sprintf("consul:8500"),
		})
		if err != nil {
			panic(err)
		}
		agent := client.Agent()

		name := "identity"
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
