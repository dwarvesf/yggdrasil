package service

import (
	"github.com/dwarvesf/yggdrasil/services/identity/service/referral"
	"github.com/dwarvesf/yggdrasil/services/identity/service/user"
)

// Service ...
type Service struct {
	UserService     user.Service
	ReferralService referral.Service
}
