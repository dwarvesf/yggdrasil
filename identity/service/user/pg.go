package user

import (
	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/yggdrasil/identity/model"
)

type pgService struct {
	db *gorm.DB
}

func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) Save(u model.User) error {
	return s.db.Create(&u).Error
}

func (s *pgService) Get(id string) (*model.User, error) {
	var u model.User
	return &u, s.db.Where("id = ?", id).First(&u).Error
}
