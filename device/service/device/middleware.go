package device

import (
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/device/model"
)

type validationMiddleware struct {
	Service
}

//ValidationMiddleware receive device service and return device service for the purpose of validate
func ValidationMiddleware() func(Service) Service {
	return func(next Service) Service {
		return &validationMiddleware{
			Service: next,
		}
	}
}

func (mw validationMiddleware) Create(d *model.Device) (err error) {
	if &d.UserID == nil {
		return ErrUserIDIsRequired
	}
	if d.DeviceID == "" {
		return ErrDeviceIDIsRequired
	}
	if !(d.Type == model.DeviceMobile || d.Type == model.DeviceWeb || d.Type == model.DeviceDesktop) {
		return ErrInvalidType
	}
	if !(d.Status == model.StatusLogout || d.Status == model.StatusOffline || d.Status == model.StatusOnline) {
		return ErrInvalidStatus
	}
	return mw.Service.Create(d)
}

func (mw validationMiddleware) Get(deviceID uuid.UUID) (*model.Device, error) {
	return mw.Service.Get(deviceID)
}

func (mw validationMiddleware) GetList(query Query) ([]model.Device, error) {
	return mw.Service.GetList(query)
}

func (mw validationMiddleware) Update(d *model.Device) (err error) {
	if d.Status > model.StatusLogout {
		return ErrInvalidStatus
	}
	return mw.Service.Update(d)
}
