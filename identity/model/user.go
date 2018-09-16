package model

import "github.com/jinzhu/gorm/dialects/postgres"

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

type User struct {
	ID          string         `json:"id" gorm:"not null"`
	Email       string         `json:"email" gorm:"default:''"`
	Username    string         `json:"username" gorm:"default:''"`
	PhoneNumber string         `json:"phone_number" gorm:"default:''"`
	LoginType   LoginType      `json:"login_type" sql:"-"`
	Password    string         `json:"password" gorm:"default:''"`
	Salt        string         `json:"salt" gorm:"default:''"`
	Status      UserStatus     `json:"status" gorm:"default:'1'"`
	Info        postgres.Jsonb `json:"info" gorm:"type:jsonb"`
}
