package service

import (
	"github.com/dwarvesf/yggdrasil/scheduler/service/message"
	"github.com/dwarvesf/yggdrasil/scheduler/service/request"
)

type Service struct {
	RequestService request.Service
	MessageService message.Service
}
