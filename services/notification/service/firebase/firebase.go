package firebase

import (
	"context"
	"strconv"

	"github.com/dwarvesf/yggdrasil/services/notification/model"

	"github.com/NaySoftware/go-fcm"
)

// FirebaseNotifier support send notify from firebase
type FirebaseNotifier struct {
	client *fcm.FcmClient
}

// Send send notify with firebase
func (f *FirebaseNotifier) Send(ctx context.Context, deviceTokens []model.DeviceToken, title, body string, data interface{}) error {
	for _, token := range deviceTokens {
		f.client.NewFcmRegIdsMsg([]string{token.Token}, data)
		f.client.SetNotificationPayload(&fcm.NotificationPayload{
			Title: title,
			Body:  body,
			Badge: strconv.Itoa(token.Badge),
		})
		f.client.Send()
	}
	return nil
}

// New new one notifier instance
func New(ctx context.Context, s string) *FirebaseNotifier {
	return &FirebaseNotifier{client: fcm.NewFcmClient(s)}
}
