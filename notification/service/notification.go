package notification

import (
	"context"

	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"golang.org/x/oauth2/google"
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
func New(ctx context.Context, credentialConfig []byte) *FirebaseNotifier {
	creds, err := google.CredentialsFromJSON(ctx, credentialConfig, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		panic(err)
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentials(creds))
	if err != nil {
		panic(err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		panic(err)
	}

	return &FirebaseNotifier{client}
}

//Support any notification provider
