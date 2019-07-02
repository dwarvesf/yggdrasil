package organization

import (
	"github.com/dwarvesf/yggdrasil/services/organization/model"
)

type Service interface {
	Create(org *model.Organization) (*model.Organization, error)
	Update(org *model.Organization) (*model.Organization, error)
	Archive(org *model.Organization) (*model.Organization, error)
	Leave(uo *model.UserOrganizations) error
	Join(uo *model.UserOrganizations) error
}
