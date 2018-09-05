// logging middleware

package add

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	Service
}

// LoggingMiddleware ...
func LoggingMiddleware(logger log.Logger) func(Service) Service {
	return func(next Service) Service {
		return &loggingMiddleware{
			logger:  log.With(logger, "service", "Post"),
			Service: next,
		}
	}
}

func (mw loggingMiddleware) Log(method string, input, output interface{}, err error, took time.Duration) error {
	return mw.logger.Log(
		"method", method,
		"input", input,
		"output", output,
		"err", err,
		"took", took,
	)
}

func (mw loggingMiddleware) Add(ctx context.Context, arg *Add) (result int, err error) {
	defer func(begin time.Time) {
		var input string
		if arg != nil {
			input = fmt.Sprintf("x=%v,y=%v", arg.X, arg.Y)
		}
		mw.Log(
			"Add",
			input,
			result,
			err,
			time.Since(begin),
		)
	}(time.Now())

	return mw.Service.Add(ctx, arg)
}
