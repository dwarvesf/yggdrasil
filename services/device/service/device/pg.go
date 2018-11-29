package device

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/device/model"
	identityModel "github.com/dwarvesf/yggdrasil/services/identity/model"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

type pgService struct {
	db *gorm.DB
}

//NewPGService create new postgres db service
func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) ValidateUser(userID uuid.UUID) error {
	err := s.db.Where("id = ?", userID).First(&identityModel.User{}).Error
	if err != nil {
		return ErrInvalidUser
	}

	return nil
}

func (s *pgService) Create(d *model.Device) error {
	return s.db.Create(&d).Error
}

func (s *pgService) Get(deviceID uuid.UUID) (*model.Device, error) {
	device := model.Device{}
	err := s.db.Where("id = ?", deviceID).First(&device).Error
	if err != nil {
		return nil, ErrDeviceNotExist
	}

	return &device, nil
}

func (s *pgService) GetList(query Query) ([]model.Device, error) {
	db := s.db

	if toolkit.IsUUIDZero(&query.UserID) {
		db = db.Where("user_id = ?", query.UserID)
	}

	var devices []model.Device
	return devices, db.Find(&devices).Error
}

func (s *pgService) Update(d *model.Device) error {
	return s.db.Model(&model.Device{}).Where("id = ?", d.ID).First(&d).Updates(&d).Error
}
