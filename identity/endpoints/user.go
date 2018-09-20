package endpoints

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/k0kubun/pp"

	"github.com/dwarvesf/yggdrasil/identity/model"
	"github.com/dwarvesf/yggdrasil/identity/service"
	"github.com/dwarvesf/yggdrasil/identity/service/user"
	util "github.com/dwarvesf/yggdrasil/identity/util"
)

//GetUserRequest ...
type GetUserRequest struct {
	ID string `json:"id"`
}

// MakeGetUserEndpoint ...
func MakeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		pp.Println(request.(GetUserRequest))
		return s.UserService.Get(&user.UserQuery{ID: request.(GetUserRequest).ID})
	}
}

//RegisterRequest ...
type RegisterRequest struct {
	LoginType   model.LoginType `json:"login_type"`
	Email       string          `json:"email,omitempty"`
	Username    string          `json:"username,omitempty"`
	PhoneNumber string          `json:"phone_number,omitempty"`
	Password    string          `json:"password"`
	Info        postgres.Jsonb  `json:"info"`
}

//RegisterResponse ...
type RegisterResponse struct {
	Token string `json:"token"`
}

//RegisterEndpoint ...
func RegisterEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)

		passwordHashed, hashError := util.HashPassword(req.Password)
		if hashError != nil {
			return nil, hashError
		}

		verifyToken := util.GenerateToken()

		user := &model.User{
			Email:       req.Email,
			Username:    req.Username,
			PhoneNumber: req.PhoneNumber,
			LoginType:   req.LoginType,
			Password:    passwordHashed,
			Token:       verifyToken,
			Info:        req.Info,
		}

		err := s.UserService.Save(user)
		if err != nil {
			return nil, err
		}

		return RegisterResponse{Token: verifyToken}, nil
	}
}

//VerifyUserRequest ...
type VerifyUserRequest struct {
	LoginType   model.LoginType `json:"login_type"`
	Email       string          `json:"email,omitempty"`
	Username    string          `json:"username,omitempty"`
	PhoneNumber string          `json:"phone_number,omitempty"`
	VerifyToken string          `json:"token"`
}

//VerifyUserResponse ...
type VerifyUserResponse struct {
	Status string `json:"status"`
}

//VerifyUserEndpoints ...
func VerifyUserEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyUserRequest)
		u := &model.User{}
		err := errors.New("")

		if req.LoginType == model.LoginTypeEmail {
			u, err = s.UserService.Get(&user.UserQuery{LoginType: model.LoginTypeEmail, Email: req.Email})
		}

		if req.LoginType == model.LoginTypePhoneNumber {
			u, err = s.UserService.Get(&user.UserQuery{LoginType: model.LoginTypePhoneNumber, PhoneNumber: req.PhoneNumber})
		}

		if req.LoginType == model.LoginTypeUsername {
			u, err = s.UserService.Get(&user.UserQuery{LoginType: model.LoginTypeUsername, Username: req.Username})
		}

		if err != nil {
			return nil, errors.New("Token invalid")
		}
		if u.Token == req.VerifyToken {
			err := s.UserService.MakeActive(u)
			if err != nil {
				return nil, errors.New("Token invalid")
			}
			return VerifyUserResponse{Status: "success"}, nil
		}
		return nil, errors.New("Token invalid")
	}
}

//LoginRequest ...
type LoginRequest struct {
	LoginType string `json:"login_type"`
	Identity  string `json:"identity"`
	Password  string `json:"password"`
}

//LoginResponse ...
type LoginResponse struct {
	JwtToken string `json:"jwt_token"`
}

//LoginEndpoint ...
func LoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)

		user, err := s.UserService.Login(req.LoginType, req.Identity)
		if err != nil {
			return nil, ErrorInvalidLogin
		}

		if user.Status != model.UserStatusActive {
			return nil, ErrorNotActive
		}

		if err := util.VerifyPassword(req.Password, user.Password); err != nil {
			return nil, ErrorInvalidLogin
		}

		loginAuth, tokenError := util.CreateAccessToken(req.Identity, string(req.LoginType))
		if tokenError != nil {
			return nil, tokenError
		}

		return LoginResponse{loginAuth}, nil
	}
}

//VerifyTokenRequest ...
type VerifyTokenRequest struct {
	JwtToken string `json:"jwt_token"`
}

//VerifyTokenResponse ...
type VerifyTokenResponse struct {
	Status string `json:"status"`
}

//VerifyTokenEndpoints ...
func VerifyTokenEndpoints(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyTokenRequest)

		if err := util.VerifyToken(req.JwtToken); err != nil {
			return nil, ErrUnauthorize
		}
		return VerifyTokenResponse{"ok"}, nil
	}
}
