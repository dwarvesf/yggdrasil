package notification

import (
	"context"

	"github.com/dwarvesf/yggdrasil/services/notification/model"
)

// Notificationer ..
type Notificationer interface {
	Send(ctx context.Context, devices []model.DeviceToken, body, title string, data interface{}) error
}
