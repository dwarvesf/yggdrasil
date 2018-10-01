package group

import (
	"github.com/dwarvesf/yggdrasil/organization/model"
)

type Service interface {
	Create(g *model.Group) (*model.Group, error)
	Update(g *model.Group) (*model.Group, error)
}
