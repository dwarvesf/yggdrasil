package organization

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

func (s *pgService) Create(org *model.Organization) (*model.Organization, error) {
	return org, s.db.Create(org).Error
}

func (s *pgService) Join(uo *model.UserOrganizations) error {
	if err := s.checkOrganization(uo.OrganizationID); err != nil {
		return err
	}

	err := s.db.Model(&model.UserOrganizations{}).
		Where("user_id = ?", uo.UserID).
		Find(&model.UserGroups{}).
		Error
	if err == nil {
		return ErrAlreadyJoined
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	return s.db.Create(uo).Error
}

func (s *pgService) Leave(uo *model.UserOrganizations) error {
	if err := s.checkOrganization(uo.OrganizationID); err != nil {
		return err
	}

	err := s.db.Model(&model.UserOrganizations{}).
		Where("user_id = ? AND organization_id = ?", uo.UserID, uo.OrganizationID).
		Updates(uo).
		Find(uo).
		Error

	if err == gorm.ErrRecordNotFound {
		return ErrHasNotJoined
	}

	return err
}

func (s *pgService) Update(org *model.Organization) (*model.Organization, error) {
	err := s.db.Model(&model.Organization{}).
		Where("id = ?", org.ID).
		Updates(org).
		Find(org).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return org, nil
}

func (s *pgService) Archive(org *model.Organization) (*model.Organization, error) {
	tx := s.db.Begin()

	// Update organization
	err := tx.Model(&model.Organization{}).
		Where("id = ?", org.ID).
		Updates(org).
		Preload("Groups").
		Find(org).
		Error
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, ErrNotFound
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update group
	groupIDs := make([]uuid.UUID, len(org.Groups))
	for _, g := range org.Groups {
		groupIDs = append(groupIDs, g.ID)
	}

	err = tx.Model(model.Group{}).
		Where("id IN (?)", groupIDs).
		Updates(model.Group{Status: model.GroupStatusInactive}).
		Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return org, nil
}

func (s *pgService) checkOrganization(id uuid.UUID) error {
	var org model.Organization

	err := s.db.Where("id = ?", id).Find(&org).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}

	if org.Status == model.OrganizationStatusInactive {
		return ErrOrganizationNotActive
	}

	return err
}
