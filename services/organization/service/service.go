package service

import (
	"github.com/dwarvesf/yggdrasil/services/organization/service/group"
	"github.com/dwarvesf/yggdrasil/services/organization/service/organization"
)

// Service ...
type Service struct {
	OrganizationService organization.Service
	GroupService        group.Service
}
