package notification

import (
	"context"
	"log"
	"os"

	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

//Notificationer ..
type Notificationer interface {
	Send() (string, error)
}

//FirebaseNotifier support send notify from firebase
type FirebaseNotifier struct {
	client *messaging.Client
}

//Send send notify with firebase
func (noti *FirebaseNotifier) Send(ctx context.Context, deviceToken, body, title string) (string, error) {
	mess := &messaging.Message{
		Data: map[string]string{
			"title": title,
			"body":  body,
		},
		Token: deviceToken,
	}
	return noti.client.Send(ctx, mess)
}

//New new one notifier instance
func New(ctx context.Context, credentialFileName string, projectID string) *FirebaseNotifier {
	opt := option.WithCredentialsFile(credentialFileName)
	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	client, createClientError := app.Messaging(ctx)
	if createClientError != nil {
		log.Fatal(createClientError)
		os.Exit(2)
	}

	return &FirebaseNotifier{client}
}

//Support any notification provider
