package user

import (
	"regexp"

	"github.com/dwarvesf/yggdrasil/identity/model"
)

// Declare Regex
const (
	emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	phoneRegex = "^[0-9]+$"
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

func (mw validationMiddleware) Save(r *model.User) (err error) {
	if r.LoginType == model.LoginTypeEmail {
		if r.Email == "" {
			return ErrEmailEmpty
		}

		//Check email format
		emailRegexp, _ := regexp.Compile(emailRegex)
		if !emailRegexp.MatchString(r.Email) {
			return ErrEmailFormat
		}

		_, err := mw.Service.Get(&UserQuery{LoginType: model.LoginTypeEmail, Email: r.Email})
		if err == nil {
			return ErrEmailExist
		}
		return mw.Service.Save(r)
	}

	if r.LoginType == model.LoginTypePhoneNumber {
		if r.PhoneNumber == "" {
			return ErrPhoneEmpty
		}

		//Check phone number format
		phoneRegexp, _ := regexp.Compile(phoneRegex)
		if !phoneRegexp.MatchString(r.PhoneNumber) {
			return ErrPhoneNumberFortmat
		}

		_, err := mw.Service.Get(&UserQuery{LoginType: model.LoginTypePhoneNumber, PhoneNumber: r.PhoneNumber})
		if err == nil {
			return ErrPhonenumberExist
		}
		return mw.Service.Save(r)
	}

	if r.Username == "" {
		return ErrUsernameEmpty
	}

	if r.Password == "" {
		return ErrPasswordEmpty
	}

	_, errX := mw.Service.Get(&UserQuery{LoginType: model.LoginTypeUsername, Username: r.Username})
	if errX == nil {
		return ErrUsernameExist
	}
	return mw.Service.Save(r)
}
