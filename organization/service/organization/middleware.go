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

func (mw validationMiddleware) Create(org *model.Organization) (*model.Organization, error) {
	if org.Name == "" {
		return nil, ErrNameEmpty
	}

	if org.Metadata == nil {
		org.Metadata = make(model.Metadata)
	}

	return mw.Service.Create(org)
}

func (mw validationMiddleware) Update(org *model.Organization) (*model.Organization, error) {
	if org.Name == "" {
		return nil, ErrNameEmpty
	}

	if org.Metadata == nil {
		org.Metadata = make(model.Metadata)
	}

	return mw.Service.Update(org)
}

func (mw validationMiddleware) Archive(org *model.Organization) (*model.Organization, error) {
	return mw.Service.Archive(org)
}

func (mw validationMiddleware) Join(uo *model.UserOrganizations) error {
	return mw.Service.Join(uo)
}

func (mw validationMiddleware) Leave(uo *model.UserOrganizations) error {
	return mw.Service.Leave(uo)
}
