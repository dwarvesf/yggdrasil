package organization

import (
	"github.com/dwarvesf/yggdrasil/organization/model"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	Create(name string) (*model.Organization, error)
	Update(orgID uuid.UUID, name string) (*model.Organization, error)
}
