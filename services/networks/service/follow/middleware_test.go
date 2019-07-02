package follow

import (
	"testing"

	"github.com/dwarvesf/yggdrasil/services/networks/model"
	"github.com/satori/go.uuid"
)

func Test_validationMiddleware_Save(t *testing.T) {
	serviceMock := &ServiceMock{
		SaveFunc: func(r *model.Follow) error {
			return nil
		},
		FindAllFunc: func(q *Query) ([]model.Follow, error) {
			return nil, nil
		},
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type args struct {
		r *model.Follow
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid test case",
			args: args{
				&model.Follow{
					FromUser: fakeFromUUID,
					ToUser:   fakeToUUID,
				},
			},
		},
		{
			name: "Miss touser",
			args: args{
				&model.Follow{
					FromUser: fakeFromUUID,
				},
			},
			wantErr: true,
		},
		{
			name: "Miss from",
			args: args{
				&model.Follow{
					ToUser: fakeToUUID,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}
			if err := mw.Save(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("ValidationMiddleware.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validationMiddleware_UnFollow(t *testing.T) {
	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	serviceMock := &ServiceMock{
		SaveFunc: func(r *model.Follow) error {
			return nil
		},
		FindAllFunc: func(q *Query) ([]model.Follow, error) {
			return nil, nil
		},
	}

	type args struct {
		fromUser uuid.UUID
		toUser   uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "From user and To user not exists",
			args: args{
				fromUser: fakeToUUID,
				toUser:   fakeFromUUID,
			},
			wantErr: true,
		},
		{
			name: "Miss touser",
			args: args{
				fromUser: fakeFromUUID,
			},
			wantErr: true,
		},
		{
			name: "Miss from",
			args: args{
				toUser: fakeToUUID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}
			if err := mw.UnFollow(tt.args.fromUser, tt.args.toUser); (err != nil) != tt.wantErr {
				t.Errorf("ValidationMiddleware.UnFollow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
