package user

import (
	"github.com/dwarvesf/yggdrasil/services/identity/model"
)

type Service interface {
	Save(r *model.User) error
	Get(userQuery *UserQuery) (*model.User, error)
	MakeActive(user *model.User) error
	Login(loginType, identity string) (*model.User, error)
}

type UserQuery struct {
	ID          string
	LoginType   model.LoginType
	Email       string
	Username    string
	PhoneNumber string
}
