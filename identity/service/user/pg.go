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

func (s *pgService) Save(u *model.User) error {
	return s.db.Create(u).Error
}

func (s *pgService) Get(userQuery *UserQuery) (*model.User, error) {
	u := &model.User{}

	db := s.db

	if userQuery.ID != "" {
		db = db.Where("id = ?", userQuery.ID)
	}

	if userQuery.Email != "" {
		db = db.Where("email = ?", userQuery.Email)
	}

	if userQuery.PhoneNumber != "" {
		db = db.Where("phone_number = ?", userQuery.PhoneNumber)
	}

	if userQuery.Username != "" {
		db = db.Where("username = ?", userQuery.Username)
	}

	if err := db.Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *pgService) MakeActive(user *model.User) error {
	err := s.db.Model(user).Updates(map[string]interface{}{"status": model.UserStatusActive, "token": ""}).Error
	return err
}

func (s *pgService) Login(loginType, identity string) (*model.User, error) {
	var u model.User

	if err := s.db.Where(map[string]interface{}{loginType: identity}).Find(&u); err.Error != nil {
		return nil, err.Error
	}

	return &u, nil
}
