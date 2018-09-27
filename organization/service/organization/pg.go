package organization

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

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

func (s *pgService) Create(name string) (*model.Organization, error) {
	org := model.Organization{
		Name: name,
	}

	err := s.db.Create(&org).Error
	if err != nil {
		return nil, err
	}

	return &org, nil
}

func (s *pgService) Update(orgID uuid.UUID, name string) (*model.Organization, error) {
	org, err := s.getOrgByID(orgID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}

		return nil, err
	}

	org.Name = name
	org.UpdatedAt = time.Now()

	return org, s.db.Save(org).Error
}

func (s *pgService) getOrgByID(orgID uuid.UUID) (*model.Organization, error) {
	var org model.Organization

	return &org, s.db.Where("id = ?", orgID).Find(&org).Error
}
