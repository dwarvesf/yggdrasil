package user

import (
	"github.com/dwarvesf/yggdrasil/identity/model"
)

type Service interface {
	Save(r model.User) error
	Get(id string) (*model.User, error)
}
