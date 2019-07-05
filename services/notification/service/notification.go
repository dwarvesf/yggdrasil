package notification

import (
	"context"
)

// Notificationer ..
type Notificationer interface {
	Send(ctx context.Context, devices []string, body, title string, data interface{}) error
}
