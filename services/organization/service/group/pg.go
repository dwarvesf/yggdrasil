package group

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/organization/model"
)

type pgService struct {
	db *gorm.DB
}

// NewPGService ...
func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) Join(ug *model.UserGroups) error {
	if err := s.checkGroup(ug.GroupID); err != nil {
		return err
	}

	err := s.db.Model(&model.UserGroups{}).
		Where("user_id = ? AND group_id = ?", ug.UserID, ug.GroupID).
		Find(&model.UserGroups{}).
		Error
	if err == nil {
		return ErrAlreadyJoined
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	return s.db.Create(ug).Error
}

func (s *pgService) Leave(ug *model.UserGroups) error {
	if err := s.checkGroup(ug.GroupID); err != nil {
		return err
	}

	err := s.db.Model(&model.UserGroups{}).
		Where("user_id = ? AND group_id = ?", ug.UserID, ug.GroupID).
		Updates(ug).
		Find(ug).
		Error

	if err == gorm.ErrRecordNotFound {
		return ErrHasNotJoined
	}

	return err
}

func (s *pgService) Create(g *model.Group) (*model.Group, error) {
	if err := s.checkOrganization(g.OrganizationID); err != nil {
		return nil, err
	}

	return g, s.db.Create(g).Error
}

func (s *pgService) Archive(g *model.Group) (*model.Group, error) {
	return s.Update(g)
}

func (s *pgService) Update(g *model.Group) (*model.Group, error) {
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

func (s *pgService) checkGroup(id uuid.UUID) error {
	var g model.Group

	err := s.db.Where("id = ?", id).Find(&g).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}

	if g.Status == model.GroupStatusInactive {
		return ErrGroupNotActive
	}

	return err
}

func (s *pgService) checkOrganization(id uuid.UUID) error {
	err := s.db.Where("id = ?", id).Find(&model.Organization{}).Error
	if err == gorm.ErrRecordNotFound {
		return ErrOrgNotFound
	}

	return err
}
