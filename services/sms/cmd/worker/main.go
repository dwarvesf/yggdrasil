package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	consul "github.com/hashicorp/consul/api"
	validator "gopkg.in/validator.v2"

	"github.com/dwarvesf/yggdrasil/logger"
	"github.com/dwarvesf/yggdrasil/services/sms/model"
	sms "github.com/dwarvesf/yggdrasil/services/sms/service"
	"github.com/dwarvesf/yggdrasil/services/sms/service/twilio"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/kafka"
)

func main() {
	svcName := "sms"
	logger := logger.NewLogger()

	logger.Info("start %v worker", svcName)
	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("consul-server:8500"),
	})
	if err != nil {
		panic(err)
	}

	if err := toolkit.RegisterService(consulClient, svcName, 0); err != nil {
		logger.Error("cannot register to consul %s", err.Error())
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

		logger.Info("send %v", svcName)
		if err := send(req.Payload, consulClient); err != nil {
			logger.Error("cannot send %v %s", svcName, err.Error())

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

func send(p model.Payload, consulClient *consul.Client) error {
	var err error
	var smsClient sms.SMS

	switch p.Provider {
	case "twilio":
		v := os.Getenv("TWILIO")
		if v == "" {
			v, err = toolkit.GetConsulValueFromKey(consulClient, p.Provider)
			if err != nil {
				return err
			}
		}

		value := model.TwilioSecret{}
		if err := json.Unmarshal([]byte(v), &value); err != nil {
			return err
		}

		smsClient = twilio.New(value.AppSid, value.AuthToken)
		return smsClient.Send(value.AppNumber, p.To, p.Content, value.AppSid)
	default:
		return errors.New("INVALID_PROVIDER")
	}
}
