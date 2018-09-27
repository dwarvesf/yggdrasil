package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"

	"github.com/dwarvesf/yggdrasil/device/db"
	"github.com/dwarvesf/yggdrasil/device/endpoints"
	serviceHttp "github.com/dwarvesf/yggdrasil/device/http"
	"github.com/dwarvesf/yggdrasil/device/middlewares"
	"github.com/dwarvesf/yggdrasil/device/service"
	"github.com/dwarvesf/yggdrasil/device/service/device"
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
			logger,
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
		logger.Log("consul", "registered", "name", name)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
