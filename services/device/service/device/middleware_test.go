package device

import (
	"net/http"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/device/model"
)

func Test_validationMiddleware_Create(t *testing.T) {
	serviceMock := &ServiceMock{
		CreateFunc: func(d *model.Device) error {
			return nil
		},
	}

	userID, errParseUUID := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7c62ab")
	if errParseUUID != nil {
		t.Errorf("failed parse uuid from string")
	}

	type args struct {
		d *model.Device
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		errorStatusCode int
	}{
		{
			name: "valid device",
			args: args{
				&model.Device{
					UserID:   userID,
					Type:     1,
					DeviceID: "42fd427e-0753-4a00-aed5-483a0c7c62ac",
					Status:   2,
				},
			},
		},
		{
			name: "invalid device caused invalid device type",
			args: args{
				&model.Device{
					UserID:   userID,
					Type:     0,
					Status:   1,
					DeviceID: "42fd427e-0753-4a00-aed5-483a0c7c62ac",
				},
			},
			wantErr:         true,
			errorStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid device caused missing user id",
			args: args{
				&model.Device{
					Type:     1,
					Status:   1,
					DeviceID: "42fd427e-0753-4a00-aed5-483a0c7c62ab",
				},
			},
			wantErr:         true,
			errorStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid device cause missing device id",
			args: args{
				&model.Device{
					UserID: userID,
					Type:   1,
					Status: 1,
				},
			},
			wantErr:         true,
			errorStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid device caused invalid status device",
			args: args{
				&model.Device{
					UserID:   userID,
					Type:     1,
					DeviceID: "42fd427e-0753-4a00-aed5-483a0c7c62ab",
					Status:   6,
				},
			},
			wantErr:         true,
			errorStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}
			err := mw.Create(tt.args.d)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("validationMiddleware.Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				status, ok := err.(interface{ StatusCode() int })
				if !ok {
					t.Errorf("validationMiddleware.Create() error %v doesn't implement StatusCode()", err)
				}
				if tt.errorStatusCode != status.StatusCode() {
					t.Errorf("validationMiddleware.Create() status = %v, want status code %v", status.StatusCode(), tt.errorStatusCode)
					return
				}

				return
			}
		})
	}
}

func Test_validationMiddleware_Update(t *testing.T) {
	serviceMock := &ServiceMock{
		UpdateFunc: func(d *model.Device) error {
			return nil
		},
	}

	type args struct {
		d *model.Device
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		errorStatusCode int
	}{
		{
			name: "valid device",
			args: args{
				&model.Device{
					Status: 3,
				},
			},
		},
		{
			name: "invalid device by error invalid device status",
			args: args{
				&model.Device{
					Status: 5,
				},
			},
			wantErr:         true,
			errorStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}

			err := mw.Update(tt.args.d)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("validationMiddleware.Update() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				status, ok := err.(interface{ StatusCode() int })
				if !ok {
					t.Errorf("validationMiddleware.Update() error %v doesn't implement StatusCode()", err)
				}
				if tt.errorStatusCode != status.StatusCode() {
					t.Errorf("validationMiddleware.Update() status = %v, want status code %v", status.StatusCode(), tt.errorStatusCode)
					return
				}

				return
			}
		})
	}
}
