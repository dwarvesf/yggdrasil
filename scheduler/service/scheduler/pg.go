package scheduler

import (
	"time"

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

func (s *pgService) DeleteRequests(ids []uint) error {
	return s.db.
		Where("id IN (?)", ids).
		Delete(model.RequestEntity{}).
		Error
}

func (s *pgService) GetRequests() ([]model.RequestEntity, error) {
	var requests []model.RequestEntity

	err := s.db.
		Where("timestamp < ?", time.Now()).
		Find(&requests).
		Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (s *pgService) SaveRequest(r model.RequestEntity) error {
	return s.db.Create(&r).Error
}
