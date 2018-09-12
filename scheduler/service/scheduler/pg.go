package scheduler

import (
	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/yggdrasil/scheduler/model"
)

type pgService struct {
	db *gorm.DB
}

func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) SaveRequest(r model.RequestEntity) error {
	return s.db.Create(&r).Error
}
