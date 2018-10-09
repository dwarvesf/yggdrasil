package friend

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/networks/model"
)

type pgService struct {
	db *gorm.DB
}

//NewPGService create service db for friend
func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) Save(o *model.Friend) error {
	return s.db.Save(o).Error
}

func (s *pgService) MakeFriend(from, to uuid.UUID) error {
	friend := model.Friend{}

	err := s.db.Where(map[string]interface{}{"from_user": from, "to_user": to}).Find(&friend).Error
	if err != nil {
		return s.Save(&model.Friend{FromUser: from, ToUser: to})
	}

	if !friend.AcceptedAt.IsZero() {
		return ErrFriendAccepted
	}

	return ErrFriendRejected
}

func (s *pgService) UnFriend(from, to uuid.UUID) error {
	friend := model.Friend{}

	err := s.db.Where(map[string]interface{}{"from_user": from, "to_user": to}).Find(&friend).Error
	if err != nil {
		return ErrRequestNotExist
	}

	return s.db.Model(&friend).Updates(map[string]interface{}{"accepted_at": time.Time{},
		"rejected_at": time.Now()}).Error
}

func (s *pgService) Accept(from, to uuid.UUID) error {
	friend := model.Friend{}

	err := s.db.Where(map[string]interface{}{"from_user": from, "to_user": to}).Find(&friend).Error
	if err != nil {
		return ErrRequestNotExist
	}

	return s.db.Model(&friend).Updates(map[string]interface{}{"accepted_at": time.Now()}).Error
}

func (s *pgService) Reject(from, to uuid.UUID) error {
	friend := model.Friend{}

	err := s.db.Where(map[string]interface{}{"from_user": from, "to_user": to}).Find(&friend).Error
	if err != nil {
		return ErrRequestNotExist
	}

	return s.db.Model(&friend).Updates(map[string]interface{}{"rejected_at": time.Now()}).Error
}

func (s *pgService) GetFriends(userID uuid.UUID) ([]model.Friend, error) {
	res := []model.Friend{}
	err := s.db.Where("from_user = ?", userID).Or("to_user = ?", userID).Not("accepted_at", time.Time{}).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
