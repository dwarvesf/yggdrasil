package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	consul "github.com/hashicorp/consul/api"

	"github.com/dwarvesf/yggdrasil/logger"
	"github.com/dwarvesf/yggdrasil/services/device/db"
	"github.com/dwarvesf/yggdrasil/services/device/endpoints"
	serviceHttp "github.com/dwarvesf/yggdrasil/services/device/http"
	"github.com/dwarvesf/yggdrasil/services/device/middlewares"
	"github.com/dwarvesf/yggdrasil/services/device/service"
	"github.com/dwarvesf/yggdrasil/services/device/service/device"
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

	pgdb, closeDB := db.New(consulClient)
	db.Migrate(pgdb)
	defer closeDB()

	var s service.Service
	{
		s = service.Service{
			DeviceService: middlewares.Compose(
				device.NewPGService(pgdb),
				device.ValidationMiddleware(),
			).(device.Service),
		}
	}
	var h http.Handler
	{
		h = serviceHttp.NewHTTPHandler(
			s,
			endpoints.MakeServerEndpoint(s),
			true,
		)
	}
	errs := make(chan error)
	go func() {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}

		agent := consulClient.Agent()

		name := "device"
		if err := agent.ServiceRegister(&consul.AgentServiceRegistration{
			Name:    name,
			Port:    port,
			Address: os.Getenv("PRIVATE_IP"),
		}); err != nil {
			panic(err)
		}
		logger.Info("registered %s", name)
	}()

	go func() {
		logger.Info("transport HTTP, addr: %s", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Info("exit", <-errs)
}
