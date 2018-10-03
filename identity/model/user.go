package model

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

type LoginType string

const (
	LoginTypeEmail       LoginType = "email"
	LoginTypeUsername    LoginType = "username"
	LoginTypePhoneNumber LoginType = "phone_number"
)

type UserStatus uint8

const (
	UserStatusInactive UserStatus = iota + 1
	UserStatusActive
	UserStatusDeactive
)

//User status 0 is inactive, 1 is active
//When login, check username, password valid, and check Status be must active
type User struct {
	ID           uuid.UUID      `json:"id" gorm:"not null"`
	Email        string         `json:"email" gorm:"default:''"`
	Username     string         `json:"username" gorm:"default:''"`
	PhoneNumber  string         `json:"phone_number" gorm:"default:''"`
	LoginType    LoginType      `json:"login_type"`
	Password     string         `json:"password" gorm:"default:''"`
	Salt         string         `json:"salt" gorm:"default:''"`
	Status       UserStatus     `json:"status" gorm:"default:'1'"`
	Info         postgres.Jsonb `json:"info" gorm:"type:jsonb"`
	Token        string         `json:"token"`
	ReferralCode string         `json:"referral_code"`
}

//BeforeCreate genergrate userID before add user
func (m *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewV4())
	return nil
}
