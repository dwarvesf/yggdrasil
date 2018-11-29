package device

import (
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/device/model"
)

//Service interface save function of device service
type Service interface {
	Create(d *model.Device) error
	ValidateUser(userID uuid.UUID) error
	Get(deviceID uuid.UUID) (*model.Device, error)
	GetList(query Query) ([]model.Device, error)
	Update(d *model.Device) error
}

//Query struct use for purpose of query data
type Query struct {
	UserID uuid.UUID
}
