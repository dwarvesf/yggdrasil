package referral

import (
	"regexp"

	"github.com/dwarvesf/yggdrasil/services/identity/model"
)

var (
	emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

type validationMiddleware struct {
	Service
}

// ValidationMiddleware ...
func ValidationMiddleware() func(Service) Service {
	return func(next Service) Service {
		return &validationMiddleware{
			Service: next,
		}
	}
}

func (mw *validationMiddleware) Save(o *model.Referral) error {
	if model.IsZero(o.FromUserID) {
		return ErrFromUserIDEmpty
	}

	if o.ToUserEmail == "" {
		return ErrToEmailEmpty
	}

	//Check email format
	emailRegexp, _ := regexp.Compile(emailRegex)
	if !emailRegexp.MatchString(o.ToUserEmail) {
		return ErrEmailFormat
	}

	return mw.Service.Save(o)
}

func (mw *validationMiddleware) DeleteReferralWithCode(code string) error {
	if code == "" {
		return ErrCodeEmpty
	}
	return mw.Service.DeleteReferralWithCode(code)
}
