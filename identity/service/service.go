package service

import (
	"github.com/dwarvesf/yggdrasil/identity/service/referral"
	"github.com/dwarvesf/yggdrasil/identity/service/user"
)

// Service ...
type Service struct {
	UserService     user.Service
	ReferralService referral.Service
}
