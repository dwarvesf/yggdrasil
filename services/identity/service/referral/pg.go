package referral

import (
	"github.com/dwarvesf/yggdrasil/services/identity/model"
	"github.com/jinzhu/gorm"
)

type pgService struct {
	db *gorm.DB
}

//NewPGService ...
func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

//Create a referral, if reffer exist will delete it before create
func (s *pgService) Save(o *model.Referral) error {
	db := s.db.Where("from_user_id = ?", o.FromUserID).Where("to_user_email = ?", o.ToUserEmail)
	res := []model.Referral{}

	if err := db.Find(&res).Error; err != nil {
		return err
	}
	if err := db.Delete(&res).Error; err != nil {
		return err
	}

	return s.db.Create(o).Error
}

//DeleteReferralWithCode when response request, record referral will be deleted
func (s *pgService) DeleteReferralWithCode(code string) error {
	return s.db.Where("code = ?", code).Delete(&model.Referral{}).Error
}

func (s *pgService) Get(q *Query) (model.Referral, error) {
	res := model.Referral{}
	return res, s.db.Where("code = ?", q.Code).First(&res).Error
}
