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
	"github.com/dwarvesf/yggdrasil/services/email/model"
	email "github.com/dwarvesf/yggdrasil/services/email/service"
	mailgun "github.com/dwarvesf/yggdrasil/services/email/service/mailgun"
	sendgrid "github.com/dwarvesf/yggdrasil/services/email/service/sendgrid"
	"github.com/dwarvesf/yggdrasil/toolkit"
	"github.com/dwarvesf/yggdrasil/toolkit/queue/nats"
)

func main() {
	svcName := "email"
	logger := logger.NewLogger()
	var err error

	logger.Info("starting email worker")
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
			logger.Info("Unable to parse request %s", err.Error())
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			logger.Error("Validate error: %s", err)
		}

		logger.Info("Sending email to %v", svcName)
		if err := sendEmail(req, vaultClient); err != nil {
			logger.Error("Cannot send email", err.Error())

			logger.Info("retry payload")
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

func sendEmail(r model.Request, vaultClient *api.Client) error {
	var emailer email.Emailer

	switch r.Payload.Provider {
	case "sendgrid":
		v := os.Getenv("SENDGRID")
		var err error
		if v == "" {
			v, err = toolkit.GetVaultValueFromKey(vaultClient, "sendgrid")
			if err != nil {
				return errors.New("can't get SENDGRID value from Vault")
			}

		}
		emailer = sendgrid.New(v)
		return emailer.Send(&r.Payload)

	case "mailgun":
		v := os.Getenv("MAILGUN")
		var err error
		if v == "" {
			v, err = toolkit.GetVaultValueFromKey(vaultClient, "mailgun")
			if err != nil {
				return errors.New("can't get MAILGUN value from Vault")
			}
		}

		value := model.MailgunSecret{}
		err = json.Unmarshal([]byte(v), &value)
		if err != nil {
			return err
		}

		emailer = mailgun.New(value.Domain, value.APIKey, value.PublicKey)
		return emailer.Send(&r.Payload)

	default:
		return errors.New("INVALID_PROVIDER")
	}
}
