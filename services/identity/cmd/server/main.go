package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	consul "github.com/hashicorp/consul/api"

	"github.com/dwarvesf/yggdrasil/logger"
	cfg "github.com/dwarvesf/yggdrasil/services/identity/cmd/config"
	"github.com/dwarvesf/yggdrasil/services/identity/db"
	"github.com/dwarvesf/yggdrasil/services/identity/endpoints"
	serviceHttp "github.com/dwarvesf/yggdrasil/services/identity/http"
	"github.com/dwarvesf/yggdrasil/services/identity/middlewares"
	"github.com/dwarvesf/yggdrasil/services/identity/service"
	"github.com/dwarvesf/yggdrasil/services/identity/service/referral"
	"github.com/dwarvesf/yggdrasil/services/identity/service/user"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

func main() {
	var (
		httpAddr = flag.String("addr", fmt.Sprintf(":%v", os.Getenv("PORT")), "HTTP listen address")
	)
	flag.Parse()

	logger := logger.NewLogger()

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
			ReferralService: middlewares.Compose(
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
		logger.Info("consul registered, name: %s", name)
	}()

	go func() {
		logger.Info("transport HTTP ,addr: %s", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Info("exit", <-errs)
}