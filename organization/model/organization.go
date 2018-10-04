package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type OrganizationStatus uint8

const (
	OrganizationStatusActive OrganizationStatus = iota + 1
	OrganizationStatusInactive
)

// Organization status 0 is inactive, 1 is active
type Organization struct {
	ID        uuid.UUID          `json:"id" gorm:"not null"`
	Name      string             `json:"name" gorm:"default:''"`
	Status    OrganizationStatus `json:"status" gorm:"default:'1'"`
	Metadata  Metadata           `json:"metadata" gorm:"type:jsonb"`
	Groups    []Group            `json:"groups"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt *time.Time         `sql:"index" json:"deleted_at"`
}

// BeforeCreate ...
func (o *Organization) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}

// BeforeSave ...
func (o *Organization) BeforeSave() error {
	o.UpdatedAt = time.Now()
	return nil
}
