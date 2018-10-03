package group

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

func (mw validationMiddleware) Create(g *model.Group) (*model.Group, error) {
	if g.Name == "" {
		return nil, ErrNameEmpty
	}

	if g.Metadata == nil {
		g.Metadata = make(model.Metadata)
	}

	return mw.Service.Create(g)
}

func (mw validationMiddleware) Update(g *model.Group) (*model.Group, error) {
	if g.Name == "" {
		return nil, ErrNameEmpty
	}

	if g.Metadata == nil {
		g.Metadata = make(model.Metadata)
	}

	return mw.Service.Update(g)
}

func (mw validationMiddleware) Join(ug *model.UserGroups) error {
	return mw.Service.Join(ug)
}

func (mw validationMiddleware) Leave(ug *model.UserGroups) error {
	return mw.Service.Leave(ug)
}
