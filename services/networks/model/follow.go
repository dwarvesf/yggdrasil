package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Follow status 0 is inactive, 1 is active
type Follow struct {
	ID        uuid.UUID `json:"id" gorm:"not null"`
	FromUser  uuid.UUID
	ToUser    uuid.UUID
	Status    uint8      `gorm:"default:'1'"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate ...
func (m *Follow) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}
