package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

//DeviceType define type of device
type DeviceType uint8

//declare device type constant
const (
	DeviceMobile DeviceType = iota + 1
	DeviceWeb
	DeviceDesktop
)

//DeviceStatus define status of device
type DeviceStatus uint8

//declare device status constant
const (
	StatusOnline DeviceStatus = iota + 1
	StatusOffline
	StatusLogout
)

//Device struct contain device model
type Device struct {
	ID        uuid.UUID      `json:"id" gorm:"not null"`
	UserID    uuid.UUID      `json:"user_id" gorm:"not null"`
	DeviceID  string         `json:"device_id" gorm:"not null"`
	Type      DeviceType     `json:"device_type"`
	Metadata  postgres.Jsonb `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
	FCMToken  string         `json:"fcm_token"`
	Status    DeviceStatus   `json:"status" gorm:"not null"`
}

//BeforeCreate genergrate userID before add user
func (m *Device) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}
