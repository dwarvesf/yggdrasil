package friend

import (
	"github.com/dwarvesf/yggdrasil/networks/model"
	uuid "github.com/satori/go.uuid"
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

//Save check userid not empty
func (mw validationMiddleware) Save(o *model.Friend) (err error) {
	if model.IsZero(o.ToUser) {
		return ErrToUserEmpty
	}

	if model.IsZero(o.FromUser) {
		return ErrFromUserEmpty
	}

	return mw.Service.Save(o)
}

//MakeFriend check userid not empty
func (mw validationMiddleware) MakeFriend(from, to uuid.UUID) (err error) {
	if model.IsZero(to) {
		return ErrToUserEmpty
	}

	if model.IsZero(from) {
		return ErrFromUserEmpty
	}

	return mw.Service.MakeFriend(from, to)
}

//UnFriend check userid not empty
func (mw validationMiddleware) UnFriend(from, to uuid.UUID) (err error) {
	if model.IsZero(to) {
		return ErrToUserEmpty
	}

	if model.IsZero(from) {
		return ErrFromUserEmpty
	}

	return mw.Service.UnFriend(from, to)
}

//Accept check userid not empty
func (mw validationMiddleware) Accept(from, to uuid.UUID) (err error) {
	if model.IsZero(to) {
		return ErrToUserEmpty
	}

	if model.IsZero(from) {
		return ErrFromUserEmpty
	}

	return mw.Service.Accept(from, to)
}

//Reject check userid not empty
func (mw validationMiddleware) Reject(from, to uuid.UUID) (err error) {
	if model.IsZero(to) {
		return ErrToUserEmpty
	}

	if model.IsZero(from) {
		return ErrFromUserEmpty
	}

	return mw.Service.Reject(from, to)
}

//GetFriend check userid not empty
func (mw validationMiddleware) GetFriends(userID uuid.UUID) (res []model.Friend, err error) {
	if model.IsZero(userID) {
		return nil, ErrUserIDInvalid
	}

	return mw.Service.GetFriends(userID)
}
