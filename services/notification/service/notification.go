package notification

import (
	"context"

	"github.com/NaySoftware/go-fcm"
)

// Notificationer ..
type Notificationer interface {
	Send(ctx context.Context, devices []string, body, title string, data interface{}) error
}

// FirebaseNotifier support send notify from firebase
type FirebaseNotifier struct {
	client *fcm.FcmClient
}

// Send send notify with firebase
func (f *FirebaseNotifier) Send(ctx context.Context, deviceTokens []string, title, body string, data interface{}) error {
	f.client.NewFcmRegIdsMsg(deviceTokens, data)
	f.client.SetNotificationPayload(&fcm.NotificationPayload{
		Title: title,
		Body:  body,
	})
	_, err := f.client.Send()
	return err
}

// New new one notifier instance
func NewFirebaseNotifier(ctx context.Context, s string) Notificationer {
	return &FirebaseNotifier{client: fcm.NewFcmClient(s)}
}
