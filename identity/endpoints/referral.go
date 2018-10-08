package endpoints

import (
	"context"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/identity/model"
	"github.com/dwarvesf/yggdrasil/identity/service"
	"github.com/dwarvesf/yggdrasil/identity/service/referral"
	util "github.com/dwarvesf/yggdrasil/identity/util"
	"github.com/go-kit/kit/endpoint"
)

//ReferralRequset send email
type ReferralRequset struct {
	FromUserID  uuid.UUID      `json:"from_user_id"`
	ToUserEmail string         `json:"to_user_email"`
	Metadata    postgres.Jsonb `json:"metadta"`
}

//ReferralReponse respone
type ReferralReponse struct {
	RefferalCode string `json:"refferal_code"`
}

//ReferralUser make referral endpoints
func ReferralUser(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ReferralRequset)
		code := util.GenerateToken()

		err := s.ReferralService.Save(&model.Referral{
			FromUserID:  req.FromUserID,
			ToUserEmail: req.ToUserEmail,
			Code:        code,
		})
		if err != nil {
			return nil, err
		}

		return ReferralReponse{code}, nil
	}
}

//ResponseReferralRequset response referral request
type ResponseReferralRequset struct {
	Code     string   `json:"code"`
	UserInfo UserInfo `json:"user_info"`
}

//ResponseReferralReponse response refferal
type ResponseReferralReponse struct {
	JwtToken string `json:"jwt_token"`
}

//UserInfo ...
type UserInfo struct {
	LoginType model.LoginType `json:"login_type"`
	Identity  string          `json:"identity"`
	Password  string          `json:"password"`
}

//ResponseReferral endpoint
func ResponseReferral(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ResponseReferralRequset)

		// check code invalid
		res, err := s.ReferralService.Get(&referral.Query{Code: req.Code})
		if err != nil {
			return nil, ErrCodeInvalid
		}

		if time.Now().Second()-res.CreatedAt.Second() > res.TTL {
			return nil, ErrTTLExpires
		}

		// Make account active
		passwordHashed, hashError := util.HashPassword(req.UserInfo.Password)
		if hashError != nil {
			return nil, hashError
		}

		user := model.User{
			LoginType: req.UserInfo.LoginType,
			Password:  passwordHashed,
			Status:    model.UserStatusActive,
		}

		switch req.UserInfo.LoginType {
		case model.LoginTypeUsername:
			user.Username = req.UserInfo.Identity
		case model.LoginTypeEmail:
			user.Email = req.UserInfo.Identity
		case model.LoginTypePhoneNumber:
			user.PhoneNumber = req.UserInfo.Identity
		}

		if err := s.UserService.Save(&user); err != nil {
			return nil, err
		}

		//delete referral
		if err := s.ReferralService.DeleteReferralWithCode(req.Code); err != nil {
			return nil, ErrDeleteRefferal
		}

		//Return jwt token
		loginAuth, tokenError := util.CreateAccessToken(req.UserInfo.Identity, string(req.UserInfo.LoginType))
		if tokenError != nil {
			return nil, tokenError
		}

		return ResponseReferralReponse{loginAuth}, nil
	}
}
