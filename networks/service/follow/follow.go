package follow

import (
	"github.com/dwarvesf/yggdrasil/networks/model"
	uuid "github.com/satori/go.uuid"
)

type Service interface {
	Save(r *model.Follow) error
	UnFollow(fromUser, toUser uuid.UUID) error
	Follow(fromUser, toUser uuid.UUID) error
	FindAll(q *Query) ([]model.Follow, error)
}

//Query handle find query
type Query struct {
	ID       uuid.UUID
	FromUser uuid.UUID
	ToUser   uuid.UUID
	Status   uint8
}
