package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// UserOrganizations ...
type UserOrganizations struct {
	ID             uuid.UUID  `json:"id" gorm:"not null"`
	OrganizationID uuid.UUID  `json:"group_id" gorm:"type:uuid REFERENCES organizations(id)"`
	UserID         uuid.UUID  `json:"user_id" gorm:"not null"`
	JoinedAt       time.Time  `json:"joined_at"`
	LeftAt         *time.Time `json:"left_at"`
	InvitedBy      *uuid.UUID `json:"invited_by"`
	InvitedAt      *time.Time `json:"invited_at"`
	KickedBy       *uuid.UUID `json:"kicked_by"`
	KickedAt       *time.Time `json:"kicked_at"`
}

// BeforeCreate ...
func (uo *UserOrganizations) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	uo.JoinedAt = time.Now()
	return nil
}
