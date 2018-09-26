package organization

import (
	"github.com/dwarvesf/yggdrasil/organization/model"
)

type Service interface {
	Save(o *model.Organization) error
}
