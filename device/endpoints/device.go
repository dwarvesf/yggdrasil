package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/device/model"
	"github.com/dwarvesf/yggdrasil/device/service"
	deviceService "github.com/dwarvesf/yggdrasil/device/service/device"
)

//CreateRequest device create request
type CreateRequest struct {
	DeviceID string             `json:"device_id"`
	UserID   uuid.UUID          `json:"user_id"`
	Type     model.DeviceType   `json:"device_type"`
	Metadata postgres.Jsonb     `json:"metadata"`
	FCMToken string             `json:"fcm_token"`
	Status   model.DeviceStatus `json:"status"`
}

//CreateResponse create device response
type CreateResponse struct {
	Status string `json:"status"`
}

//MakeCreateDeviceEndpoint create device endpoint
func MakeCreateDeviceEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			req    = request.(CreateRequest)
			device = &model.Device{
				DeviceID: req.DeviceID,
				Type:     req.Type,
				Metadata: req.Metadata,
				UserID:   req.UserID,
				FCMToken: req.FCMToken,
				Status:   req.Status,
			}
		)

		err := s.DeviceService.ValidateUser(device.UserID)
		if err != nil {
			return nil, err
		}

		err = s.DeviceService.Create(device)
		if err != nil {
			return nil, err
		}

		return CreateResponse{Status: "success"}, nil
	}
}

//GetRequest prepare request data for get device by device id
type GetRequest struct {
	DeviceID uuid.UUID
}

//MakeGetDeviceEndpoint create endpoint for get device by device id
func MakeGetDeviceEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			req      = request.(GetRequest)
			deviceID = req.DeviceID
		)

		return s.DeviceService.Get(deviceID)
	}
}

//GetListDeviceRequest prepare request for get list device
type GetListDeviceRequest struct {
	UserID uuid.UUID
}

//GetListDeviceResponse create response for get list device
type GetListDeviceResponse struct {
	Devices []model.Device `json:"device"`
}

//MakeGetListDeviceEndpoint create endpoint for get device
func MakeGetListDeviceEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			req    = request.(GetListDeviceRequest)
			userID = req.UserID
			query  = deviceService.Query{UserID: userID}
		)

		err := s.DeviceService.ValidateUser(userID)
		if err != nil {
			return nil, err
		}

		devices, err := s.DeviceService.GetList(query)
		if err != nil {
			return nil, err
		}

		return GetListDeviceResponse{Devices: devices}, nil
	}
}

//UpdateRequest create request data for update device
type UpdateRequest struct {
	ID       uuid.UUID          `json:"-"`
	Status   model.DeviceStatus `json:"status"`
	Metadata postgres.Jsonb     `json:"metadata"`
	FCMToken string             `json:"fcm_token"`
}

//UpdateResponse create response for update device
type UpdateResponse struct {
	Status string `json:"status"`
}

//MakeUpdateDeviceEndpoint create endpoint for update device
func MakeUpdateDeviceEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var (
			req       = request.(UpdateRequest)
			newDevice = &model.Device{
				ID:       req.ID,
				Status:   req.Status,
				Metadata: req.Metadata,
				FCMToken: req.FCMToken,
			}
		)

		err := s.DeviceService.Update(newDevice)
		if err != nil {
			return nil, err
		}

		return UpdateResponse{Status: "success"}, nil
	}
}
