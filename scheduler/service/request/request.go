package request

import (
	"github.com/dwarvesf/yggdrasil/scheduler/model"
)

type Service interface {
	SaveRequest(r model.RequestEntity) error
	GetRequests() ([]model.RequestEntity, error)
	DeleteRequests(ids []uint) error
}
