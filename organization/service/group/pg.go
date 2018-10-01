package group

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

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

func (s *pgService) Create(g *model.Group) (*model.Group, error) {
	if err := s.checkOrganization(g.OrganizationID); err != nil {
		return nil, err
	}

	return g, s.db.Create(g).Error
}

func (s *pgService) Update(g *model.Group) (*model.Group, error) {
	if err := s.checkOrganization(g.OrganizationID); err != nil {
		return nil, err
	}

	err := s.db.Model(&model.Group{}).
		Where("id = ?", g.ID).
		Updates(g).
		Find(g).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return g, nil
}

func (s *pgService) checkOrganization(id uuid.UUID) error {
	err := s.db.Where("id = ?", id).Find(&model.Organization{}).Error
	if err == gorm.ErrRecordNotFound {
		return ErrOrgNotFound
	}

	return err
}
