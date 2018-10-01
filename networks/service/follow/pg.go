package follow

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/networks/model"
)

type pgService struct {
	db *gorm.DB
}

func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) Save(o *model.Follow) error {
	return s.db.Create(o).Error
}

func (s *pgService) UnFollow(fromUser, toUser uuid.UUID) error {
	follower := model.Follow{}

	err := s.db.Where(map[string]interface{}{"from_user": fromUser, "to_user": toUser}).Find(&follower).Error
	if err != nil {
		return err
	}

	return s.db.Model(follower).Update("status", 0).Error
}

func (s *pgService) Follow(fromUser, toUser uuid.UUID) error {
	follower := model.Follow{}

	err := s.db.Where(map[string]interface{}{"from_user": fromUser, "to_user": toUser}).Find(&follower).Error
	if err != nil {
		return s.Save(&model.Follow{FromUser: fromUser, ToUser: toUser})
	}

	return s.db.Model(follower).Update("status", 1).Error
}

func (s *pgService) FindAll(q *Query) ([]model.Follow, error) {
	res := []model.Follow{}
	db := s.db

	if !model.IsZero(q.ID) {
		db = db.Where("id = ?", q.ID)
	}
	if !model.IsZero(q.FromUser) {
		db = db.Where("from_user = ?", q.FromUser)
	}
	if !model.IsZero(q.ToUser) {
		db = db.Where("to_user = ?", q.ToUser)
	}
	if q.Status != 0 {
		db = db.Where("status = ?", q.Status)
	}

	err := db.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
