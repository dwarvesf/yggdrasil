package model

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	uuid "github.com/satori/go.uuid"
)

type Referral struct {
	ID          uuid.UUID      `json:"id" gorm:"not null"`
	TTL         time.Time      `json:"ttl"`
	FromUserID  uuid.UUID      `json:"from_user_id"`
	Code        string         `json:"code"`
	ToUserEmail string         `json:"to_user_email"`
	Metadata    postgres.Jsonb `json:"metadata"`
}
