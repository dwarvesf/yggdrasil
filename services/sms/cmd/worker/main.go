package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/signal"

	"github.com/hashicorp/vault/api"
	"github.com/nats-io/stan.go"
	"gopkg.in/go-playground/validator.v9"

	"github.com/dwarvesf/yggdrasil/logger"
	"github.com/dwarvesf/yggdrasil/services/sms/model"
	sms "github.com/dwarvesf/yggdrasil/services/sms/service"
	"github.com/dwarvesf/yggdrasil/services/sms/service/twilio"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/nats"
)

func main() {
	svcName := "sms"
	logger := logger.NewLogger()
	var err error

	logger.Info("Start %v worker", svcName)
	vaultClient, err := toolkit.NewVaultClient()
	if err != nil {
		logger.Error("Cannot connect to Vault %s", err.Error())
		panic(err)
	}

	clientID := os.Getenv("CLIENT_ID")
	cluster := os.Getenv("NATS_CLUSTER")
	if cluster == "" {
		cluster, err = toolkit.GetVaultValueFromKey(vaultClient, "nats_cluster")
		if err != nil {
			logger.Error("Cannot url from Vault %s", err.Error())
			panic(err)
		}
	}
	url := os.Getenv("NATS_URL")
	if url == "" {
		url, err = toolkit.GetVaultValueFromKey(vaultClient, "nats_url")
		if err != nil {
			logger.Error("Cannot url from Vault %s", err.Error())
			panic(err)
		}
	}

	q := nats.New(cluster, clientID, url)
	r, err := q.NewReader(svcName)
	if err != nil {
		logger.Error("Cannot connect to nats-streaming %s", err.Error())
		panic(err)
	}
	defer r.Close()

	w, err := q.NewWriter()
	if err != nil {
		logger.Error("Cannot connect to nats-streaming %s", err.Error())
		panic(err)
	}
	defer w.Close()

	handleMsg := func(msg *stan.Msg) {
		logger.Info("Received request %v", string(msg.Data))

		var req model.Request
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			logger.Info("Cannot parse request %s", err.Error())
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			logger.Error("Validate error: %s", err)
		}

		logger.Info("Sending message to %v", svcName)
		if err := send(req.Payload, vaultClient); err != nil {
			logger.Error("Cannot send %v %s", svcName, err.Error())

			logger.Info("Create retry payload")
			message, err := toolkit.CreateRetryMessage(svcName, req.Payload, req.Retry)
			if err != nil {
				logger.Error("Cannot create retry payload %s", err.Error())
			}
			w.Write("scheduler", message)
		}

		msg.Ack()
	}

	if err := r.Read(handleMsg); err != nil {
		logger.Error("Cannot read from nats-streaming %s", err.Error())
		panic(err)
	}

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			logger.Info("Received an interrupt, closing connection...\n\n")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func send(p model.Payload, vaultClient *api.Client) error {
	var smsClient sms.SMS

	switch p.Provider {
	case "twilio":
		var err error
		v := os.Getenv("TWILIO")
		if v == "" {
			v, err = toolkit.GetVaultValueFromKey(vaultClient, "twilio")
			if err != nil {
				return errors.New("can't get TWILIO value from Vault")
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
