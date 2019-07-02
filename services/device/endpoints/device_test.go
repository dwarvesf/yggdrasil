package endpoints

import (
	"context"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/device/model"
	"github.com/dwarvesf/yggdrasil/services/device/service"
	deviceService "github.com/dwarvesf/yggdrasil/services/device/service/device"
)

func TestMakeCreateDeviceEndpoint(t *testing.T) {
	mock := service.Service{
		DeviceService: &deviceService.ServiceMock{
			CreateFunc: func(d *model.Device) error {
				return nil
			},
			ValidateUserFunc: func(userID uuid.UUID) error {
				return nil
			},
		},
	}
	userID, errParseUUID := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7c62ab")
	if errParseUUID != nil {
		t.Errorf("failed parse uuid from string")
	}
	type args struct {
		req CreateRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{

			name: "create device endpoint success",
			args: args{
				CreateRequest{
					UserID:   userID,
					DeviceID: "bf2617d3-8396-4419-9fcd-94e7727f0f93",
					Type:     1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFunc := MakeCreateDeviceEndpoint(mock)
			_, err := gotFunc(context.Background(), tt.args.req)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMakeGetDeviceEndpoint(t *testing.T) {
	mock := service.Service{
		DeviceService: &deviceService.ServiceMock{
			GetFunc: func(deviceID uuid.UUID) (*model.Device, error) {
				device := model.Device{ID: deviceID}
				return &device, nil
			},
		},
	}
	deviceID, errParseUUID := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7c62ab")
	if errParseUUID != nil {
		t.Errorf("failed parse uuid from string")
	}
	type args struct {
		req GetRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success create get device endpoint",
			args: args{
				GetRequest{
					DeviceID: deviceID,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFunc := MakeGetDeviceEndpoint(mock)
			_, err := gotFunc(context.Background(), tt.args.req)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMakeGetListDeviceEndpoint(t *testing.T) {
	mock := service.Service{
		DeviceService: &deviceService.ServiceMock{
			GetListFunc: func(query deviceService.Query) ([]model.Device, error) {
				device := []model.Device{}
				return device, nil
			},
			ValidateUserFunc: func(userID uuid.UUID) error {
				return nil
			},
		},
	}
	userID, errParseUUID := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7c62ab")
	if errParseUUID != nil {
		t.Errorf("failed parse uuid from string")
	}
	type args struct {
		req GetListDeviceRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success create get device list endpoint",
			args: args{
				GetListDeviceRequest{
					UserID: userID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			gotFunc := MakeGetListDeviceEndpoint(mock)
			_, err := gotFunc(context.Background(), tt.args.req)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMakeUpdateDeviceEndpoint(t *testing.T) {
	mock := service.Service{
		DeviceService: &deviceService.ServiceMock{
			UpdateFunc: func(d *model.Device) error {
				return nil
			},
		},
	}
	type args struct {
		req UpdateRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				UpdateRequest{
					Status:   1,
					FCMToken: "token",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			gotFunc := MakeUpdateDeviceEndpoint(mock)
			_, err := gotFunc(context.Background(), tt.args.req)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
