package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/yggdrasil/services/device/service"
)

//Endpoints contain endpoint for device service
type Endpoints struct {
	CreateDevice  endpoint.Endpoint
	GetDevice     endpoint.Endpoint
	GetListDevice endpoint.Endpoint
	UpdateDevice  endpoint.Endpoint
}

//MakeServerEndpoint create server endpoint for device service
func MakeServerEndpoint(s service.Service) Endpoints {
	return Endpoints{
		CreateDevice:  MakeCreateDeviceEndpoint(s),
		GetDevice:     MakeGetDeviceEndpoint(s),
		GetListDevice: MakeGetListDeviceEndpoint(s),
		UpdateDevice:  MakeUpdateDeviceEndpoint(s),
	}
}
