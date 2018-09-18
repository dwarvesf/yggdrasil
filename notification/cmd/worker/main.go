package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
	validator "gopkg.in/validator.v2"

	"github.com/dwarvesf/yggdrasil/notification/model"
	notification "github.com/dwarvesf/yggdrasil/notification/service"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"
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
		logger.Log("server", "notification")
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
		name := "notification"
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
		logger.Log("consul", "registering", "name", name)

		if err := toolkit.RegisterService(consulClient, name, port); err != nil {
			panic(err)
		}
	}()

	go func() {
		var q queue.Queue
		q = kafka.New(consulClient)
		r := q.NewReader("notification")
		defer r.Close()

		for {
			b, err := r.Read()
			if err != nil {
				logger.Log("error", err.Error())
				continue
			}

			// TODO: simplify main function
			var req model.Request
			if err = json.Unmarshal(b, &req); err != nil {
				logger.Log("error", err.Error())
				continue
			}
			if err := validator.Validate; err != nil {
				logger.Log("error", err)
				continue
			}

			ctx := context.Background()

			switch req.Provider {
			case "firebase":
				projectID := os.Getenv("PROJECT_ID")
				if projectID == "" {
					var projectIDErr error
					projectID, projectIDErr = toolkit.GetConsulValueFromKey(consulClient, "project_id")
					if projectIDErr != nil {
						logger.Log("exit", projectIDErr)
						os.Exit(3)
					}
				}
				firebaseNotifier := notification.New(ctx, os.Getenv("CREDENTIAL_FILE"), projectID)

				res, sendErr := firebaseNotifier.Send(ctx, req.DeviceToken, req.Body, req.Title)
				if sendErr != nil {
					logger.Log(sendErr)
					continue
				}
				logger.Log(res)

				//Case use different notification provider, must send notify depend on req.DeviceType
			default:
				logger.Log("Provider not support")
			}
		}
	}()

	logger.Log("exit", <-errs)
}
