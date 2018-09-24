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

func (s *pgService) Save(o *model.Organization) error {
	return s.db.Create(o).Error
}
