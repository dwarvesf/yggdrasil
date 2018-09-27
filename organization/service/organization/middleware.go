package organization

import (
	"github.com/dwarvesf/yggdrasil/organization/model"
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

func (mw validationMiddleware) Save(o *model.Organization) error {
	if o.Name == "" {
		return ErrNameEmpty
	}

	return mw.Service.Save(o)
}
