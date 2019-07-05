package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	consul "github.com/hashicorp/consul/api"
	validator "gopkg.in/validator.v2"

	"github.com/dwarvesf/yggdrasil/logger"
	"github.com/dwarvesf/yggdrasil/services/notification/model"
	notification "github.com/dwarvesf/yggdrasil/services/notification/service"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"
)

func main() {
	logger := logger.NewLogger()
	logger.Info("start notification worker")

	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("consul-server:8500"),
	})
	if err != nil {
		panic(err)
	}

	svcName := "notification"
	if err := toolkit.RegisterService(consulClient, svcName, 0); err != nil {
		panic(err)
	}

	var q queue.Queue
	q = kafka.New(consulClient)
	r := q.NewReader(svcName)
	w := q.NewWriter("scheduler")
	defer r.Close()
	defer w.Close()
	for {
		b, err := r.Read()
		if err != nil {
			logger.Error("cannot read from kafka %s", err.Error())
			continue
		}

		logger.Info("received request %v", string(b))
		var req model.Request
		if err = json.Unmarshal(b, &req); err != nil {
			logger.Info("cannot parse request %s", err.Error())
			continue
		}
		if err := validator.Validate(req); err != nil {
			logger.Error("validate error: %s", err)
			continue
		}

		logger.Info("send notification")
		err = sendNotification(req.Payload, consulClient)
		if err != nil {
			logger.Error("cannot send notification %s", err.Error())

			logger.Info("create retry payload")
			message, err := toolkit.CreateRetryMessage(svcName, req.Payload, req.Retry)
			if err != nil {
				logger.Error("cannot create retry payload %s", err.Error())
				continue
			}
			w.Write(svcName, message)
		}
	}
}

func sendNotification(p model.Payload, consulClient *consul.Client) error {
	ctx := context.Background()

	//get credentials
	v := os.Getenv("NOTIFICATION")
	if v == "" {
		v, _ = toolkit.GetConsulValueFromKey(consulClient, "notification")
	}

	var n notification.Notificationer
	switch p.Provider {
	case "firebase":
		n = notification.NewFirebaseNotifier(ctx, v)
		return n.Send(ctx, []string{p.DeviceToken}, p.Title, p.Body, p.Data)
	default:
		return errors.New("INVALID_PROVIDER")
	}
}
