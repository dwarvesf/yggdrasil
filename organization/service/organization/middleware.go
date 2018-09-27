package organization

import (
	"github.com/dwarvesf/yggdrasil/organization/model"
	uuid "github.com/satori/go.uuid"
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

func (mw validationMiddleware) Create(name string) (*model.Organization, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}

	return mw.Service.Create(name)
}

func (mw validationMiddleware) Update(orgID uuid.UUID, name string) (*model.Organization, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}

	return mw.Service.Update(orgID, name)
}
