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
	"github.com/dwarvesf/yggdrasil/services/organization/db"
	serviceHttp "github.com/dwarvesf/yggdrasil/services/organization/http"
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

	h := serviceHttp.NewHTTPHandler(pgdb, true)

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
		logger.Info("consul registered, name: %s", name)
	}()

	go func() {
		logger.Info("transport HTTP ,addr: %s", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Info("exit", <-errs)
}
