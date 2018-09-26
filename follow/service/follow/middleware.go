package follow

import (
	"github.com/dwarvesf/yggdrasil/follow/model"
)

type validationMiddleware struct {
	Service
}

// ValidationMiddleware ...
func ValidationMiddleware() func(Service) Service {
	return func(next Service) Service {
		return &validationMiddleware{
			Service: next,
		}
	}
}

func (mw validationMiddleware) Save(o *model.Follow) (err error) {
	return mw.Service.Save(o)
}
