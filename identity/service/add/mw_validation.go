package add

import (
	"context"
)

type validationMiddleware struct {
	Service
}

func ValidationMiddleware() func(Service) Service {
	return func(next Service) Service {
		return &validationMiddleware{
			Service: next,
		}
	}
}

func (mw validationMiddleware) Add(ctx context.Context, arg *Add) (result int, err error) {
	if arg == nil {
		return -1, ErrInvalidParams
	}
	return mw.Service.Add(ctx, arg)
}
