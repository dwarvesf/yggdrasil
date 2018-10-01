package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type GroupStatus uint8

const (
	GroupStatusInactive GroupStatus = iota + 1
	GroupStatusActive
)

// Group status 0 is inactive, 1 is active
type Group struct {
	ID             uuid.UUID    `json:"id" gorm:"not null"`
	Name           string       `json:"username" gorm:"default:''"`
	Status         GroupStatus  `json:"status" gorm:"default:'1'"`
	Organization   Organization `json:"organization"`
	OrganizationID uuid.UUID    `json:"organization_id"`
	Metadata       Metadata     `json:"metadata" gorm:"type:jsonb"`
	CreatedBy      uuid.UUID    `json:"created_by" gorm:"not null"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `sql:"index" json:"deleted_at"`
}

// BeforeCreate ...
func (m *Group) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}
