package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	consul "github.com/hashicorp/consul/api"
	validator "gopkg.in/validator.v2"

	"github.com/dwarvesf/yggdrasil/logger"
	cfg "github.com/dwarvesf/yggdrasil/services/notification/cmd/config"
	"github.com/dwarvesf/yggdrasil/services/notification/model"
	notification "github.com/dwarvesf/yggdrasil/services/notification/service"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"
)

func main() {
	logger := logger.NewLogger()

	errs := make(chan error)
	go func() {
		logger.Info("starting notification")
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

	//get fcm credentials
	fcmCredentials, err := toolkit.GetConsulValueFromKey(consulClient, "fcm_credentials")
	sDec, _ := b64.StdEncoding.DecodeString(fcmCredentials)
	cfg.FirebaseCredentials = sDec

	svcName := "notification"
	go func() {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
		logger.Info("consul registering, name: %s", svcName)

		if err := toolkit.RegisterService(consulClient, svcName, port); err != nil {
			panic(err)
		}
	}()

	go func() {
		var q queue.Queue
		q = kafka.New(consulClient)
		r := q.NewReader(svcName)
		w := q.NewWriter("scheduler")
		defer r.Close()
		defer w.Close()

		for {
			b, err := r.Read()
			if err != nil {
				logger.Error("unable to read from kafka %s", err.Error())
				continue
			}

			var req model.Request
			if err = json.Unmarshal(b, &req); err != nil {
				logger.Info("unable to parse request %s", err.Error())
				continue
			}
			if err := validator.Validate; err != nil {
				logger.Error("Validator error: %s", err)
				continue
			}

			logger.Info("sending email")
			_, err = sendNotification(req.Payload, consulClient)
			if err != nil {
				message, err := toolkit.CreateRetryMessage("notification", req.Payload, req.Retry)
				if err != nil {
					logger.Error("unable to send an notification %s", err.Error())
					continue
				}

				w.Write("notification", message)
				logger.Info("retry payload")
			}
		}
	}()

	logger.Error("exit", <-errs)
}

func sendNotification(p model.Payload, consulClient *consul.Client) (string, error) {
	ctx := context.Background()

	switch p.Provider {
	case "firebase":
		firebaseNotifier := notification.New(ctx, cfg.FirebaseCredentials)
		return firebaseNotifier.Send(ctx, p.DeviceToken, p.Body, p.Title)
	default:
		return "", errors.New("INVALID_PROVIDER")
	}
}
