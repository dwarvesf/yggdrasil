package service

import (
	"github.com/dwarvesf/yggdrasil/scheduler/service/request"
	"github.com/dwarvesf/yggdrasil/scheduler/service/stream"
)

type Service struct {
	RequestService request.Service
	StreamService  stream.Service
}
