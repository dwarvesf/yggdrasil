package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"

	"github.com/dwarvesf/yggdrasil/scheduler/db"
	"github.com/dwarvesf/yggdrasil/scheduler/scheduler"
	"github.com/dwarvesf/yggdrasil/scheduler/service"
	"github.com/dwarvesf/yggdrasil/scheduler/service/message"
	"github.com/dwarvesf/yggdrasil/scheduler/service/request"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	errs := make(chan error)
	go func() {
		logger.Log("worker", "scheduler")
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("consul:8500"),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		name := "scheduler"
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
		logger.Log("consul", "registering", "name", name)

		if err := toolkit.RegisterService(consulClient, name, port); err != nil {
			panic(err)
		}
	}()

	pgdb, closeDB := db.New(consulClient)
	db.Migrate(pgdb)
	s := service.Service{
		RequestService: request.NewPGService(pgdb),
		MessageService: message.NewKafkaService(consulClient),
	}
	defer closeDB()

	sch := scheduler.NewScheduler(s, logger)
	go sch.HandleRequests(2 * time.Minute)
	go sch.ListenMessages()

	logger.Log("exit", <-errs)
}
