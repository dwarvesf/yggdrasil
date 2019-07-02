package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

//Referral storage invition
type Referral struct {
	ID          uuid.UUID      `json:"id" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty" gorm:"default null"`
	TTL         int            `json:"ttl" gorm:"not null"`
	FromUserID  uuid.UUID      `json:"from_user_id" gorm:"not null"`
	Code        string         `json:"code" gorm:"not null"`
	ToUserEmail string         `json:"to_user_email" gorm:"not null"`
	Metadata    postgres.Jsonb `json:"metadata" gorm:"type:jsonb"`
}

//BeforeCreate genergrate userID before add user
func (m *Referral) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	// TTL is 1 day
	scope.SetColumn("TTL", 3600*24)
	return nil
}
