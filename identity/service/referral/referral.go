package referral

import (
	"time"

	"github.com/dwarvesf/yggdrasil/identity/model"
)

//Service ...
type Service interface {
	Save(o *model.Referral) error
	DeleteReferralWithCode(code string) error
	Get(q *Query) (model.Referral, error)
}

//Query ...
type Query struct {
	Code string
	TTL  time.Time
}
