package service

import (
	"github.com/dwarvesf/yggdrasil/organization/service/group"
	"github.com/dwarvesf/yggdrasil/organization/service/organization"
)

// Service ...
type Service struct {
	OrganizationService organization.Service
	GroupService        group.Service
}
