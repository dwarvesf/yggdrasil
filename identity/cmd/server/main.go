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

	cfg "github.com/dwarvesf/yggdrasil/identity/cmd/config"
	"github.com/dwarvesf/yggdrasil/identity/db"
	"github.com/dwarvesf/yggdrasil/identity/endpoints"
	serviceHttp "github.com/dwarvesf/yggdrasil/identity/http"
	"github.com/dwarvesf/yggdrasil/identity/middlewares"
	"github.com/dwarvesf/yggdrasil/identity/service"
	"github.com/dwarvesf/yggdrasil/identity/service/referral"
	"github.com/dwarvesf/yggdrasil/identity/service/user"
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

	// FIXME: replace this with `postgres.New()`
	pgdb, closeDB := db.New(consulClient)
	db.Migrate(pgdb)
	defer closeDB()

	var s service.Service
	{
		s = service.Service{
			UserService: middlewares.Compose(
				user.NewPGService(pgdb),
				user.ValidationMiddleware(),
			).(user.Service),
			ReferrService: middlewares.Compose(
				referral.NewPGService(pgdb),
				referral.ValidationMiddleware(),
			).(referral.Service),
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
