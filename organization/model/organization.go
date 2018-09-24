package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

type OrganizationStatus uint8

const (
	OrganizationStatusInactive OrganizationStatus = iota + 1
	OrganizationStatusActive
)

// Organization status 0 is inactive, 1 is active
type Organization struct {
	ID        uuid.UUID          `json:"id" gorm:"not null"`
	Name      string             `json:"username" gorm:"default:''"`
	Status    OrganizationStatus `json:"status" gorm:"default:'1'"`
	Metadata  postgres.Jsonb     `json:"metadata" gorm:"type:jsonb"`
	Groups    []Group            `json:"organizations"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt *time.Time         `sql:"index" json:"deleted_at"`
}

// BeforeCreate ...
func (m *Organization) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}
