package service

import (
	"github.com/dwarvesf/yggdrasil/services/networks/service/follow"
	"github.com/dwarvesf/yggdrasil/services/networks/service/friend"
)

// Service ...
type Service struct {
	FollowService follow.Service
	FriendService friend.Service
}
