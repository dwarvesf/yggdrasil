package follow

import (
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/follow/model"
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

func (mw validationMiddleware) Save(r *model.Follow) (err error) {
	if model.IsZero(r.FromUser) {
		return ErrFromUserUUIDEmpty
	}
	if model.IsZero(r.ToUser) {
		return ErrToUserUUIDEmpty
	}

	//Check From user and ToUser is unique
	res, _ := mw.Service.FindAll(&Query{
		FromUser: r.FromUser,
		ToUser:   r.ToUser,
	})

	if len(res) != 0 {
		return ErrorCreateFollow
	}

	return mw.Service.Save(r)
}

func (mw validationMiddleware) UnFollow(fromUser, toUser uuid.UUID) (err error) {
	if model.IsZero(fromUser) {
		return ErrFromUserUUIDEmpty
	}
	if model.IsZero(toUser) {
		return ErrToUserUUIDEmpty
	}

	//Check From user and ToUser is unique
	res, _ := mw.Service.FindAll(&Query{
		FromUser: fromUser,
		ToUser:   toUser,
	})

	if len(res) == 0 {
		return ErrUnfollow
	}

	return mw.Service.UnFollow(fromUser, toUser)
}

func (mw validationMiddleware) Follow(fromUser, toUser uuid.UUID) (err error) {
	if model.IsZero(fromUser) {
		return ErrFromUserUUIDEmpty
	}
	if model.IsZero(toUser) {
		return ErrToUserUUIDEmpty
	}

	return mw.Service.Follow(fromUser, toUser)
}
