package organization

import (
	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/yggdrasil/organization/model"
)

type pgService struct {
	db *gorm.DB
}

func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) Create(org *model.Organization) (*model.Organization, error) {
	return org, s.db.Create(org).Error
}

func (s *pgService) Update(org *model.Organization) (*model.Organization, error) {
	err := s.db.Model(&model.Organization{}).
		Where("id = ?", org.ID).
		Updates(org).
		Find(org).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return org, nil
}
