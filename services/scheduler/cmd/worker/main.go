package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	consul "github.com/hashicorp/consul/api"

	"github.com/dwarvesf/yggdrasil/logger"
	"github.com/dwarvesf/yggdrasil/services/scheduler/db"
	"github.com/dwarvesf/yggdrasil/services/scheduler/scheduler"
	"github.com/dwarvesf/yggdrasil/services/scheduler/service"
	"github.com/dwarvesf/yggdrasil/services/scheduler/service/request"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"
)

func main() {
	logger := logger.NewLogger()

	errs := make(chan error)
	go func() {
		logger.Info("starting scheduler worker")
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("consul-server:8500"),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		name := "scheduler"
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			logger.Error("unable to get port %s", err.Error())
			panic(err)
		}
		logger.Info("registering %s to consul", name)

		if err := toolkit.RegisterService(consulClient, name, port); err != nil {
			panic(err)
		}
	}()

	pgdb, closeDB := db.New(consulClient)
	db.Migrate(pgdb)
	s := service.Service{
		RequestService: request.NewPGService(pgdb),
		QueueService:   kafka.New(consulClient),
	}
	defer closeDB()

	sch := scheduler.NewScheduler(s, *logger)
	go sch.HandleRequests(2 * time.Minute)
	go sch.ListenMessages()

	logger.Error("exit", <-errs)
}
