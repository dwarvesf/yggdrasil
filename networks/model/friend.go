package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

//Friend ...
type Friend struct {
	ID          uuid.UUID  `json:"id" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	FromUser    uuid.UUID  `json:"from" gorm:"not null"`
	ToUser      uuid.UUID  `json:"to" gorm:"not null"`
	RequestedAt time.Time  `json:"requested_at"`
	AcceptedAt  time.Time  `json:"accepted_at"`
	RejectedAt  time.Time  `json:"rejected_at"`
}

//BeforeCreate ...
func (f *Friend) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	scope.SetColumn("RequestedAt", time.Now())
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}
