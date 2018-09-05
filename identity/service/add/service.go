package add

import (
	"context"
)

type Service interface {
	Add(ctx context.Context, arg *Add) (int, error)
}
