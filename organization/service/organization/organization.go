package organization

import (
	"github.com/dwarvesf/yggdrasil/organization/model"
)

type Service interface {
	Save(r *model.Organization) error
}
