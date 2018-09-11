package scheduler

import "github.com/jinzhu/gorm"

type pgService struct {
	db *gorm.DB
}

func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

// func (s *pgService) Send() error
