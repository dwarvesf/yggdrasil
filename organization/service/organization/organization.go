package organization

import (
	"github.com/dwarvesf/yggdrasil/organization/model"
)

type Service interface {
	Create(org *model.Organization) (*model.Organization, error)
	Update(org *model.Organization) (*model.Organization, error)
}
