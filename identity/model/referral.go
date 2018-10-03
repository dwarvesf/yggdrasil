package model

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type Referral struct {
	ID          uuid.UUID      `json:"id" gorm:"not null"`
	TTL         time.Time      `json:"ttl" gorm:"not null"`
	FromUserID  uuid.UUID      `json:"from_user_id" gorm:"not null"`
	Code        string         `json:"code" gorm:"not null"`
	ToUserEmail string         `json:"to_user_email" gorm:"not null"`
	Metadata    postgres.Jsonb `json:"metadata"`
}
