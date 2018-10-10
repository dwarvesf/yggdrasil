package friend

import (
	"github.com/dwarvesf/yggdrasil/networks/model"
	uuid "github.com/satori/go.uuid"
)

//Service storage friend service
type Service interface {
	Save(o *model.Friend) error
	MakeFriend(from, to uuid.UUID) error
	UnFriend(from, to uuid.UUID) error
	Accept(from, to uuid.UUID) error
	Reject(from, to uuid.UUID) error
	GetFriends(userID uuid.UUID) ([]model.Friend, error)
}
