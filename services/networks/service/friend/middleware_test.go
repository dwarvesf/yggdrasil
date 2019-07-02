package friend

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/networks/model"
)

func Test_validationMiddleware_Save(t *testing.T) {
	serviceMock := &ServiceMock{
		SaveFunc: func(o *model.Friend) error {
			return nil
		},
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type args struct {
		o *model.Friend
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid testcase",
			args: args{
				&model.Friend{
					FromUser: fakeFromUUID,
					ToUser:   fakeToUUID,
				},
			},
		},

		{
			name: "FromUser UUID empty",
			args: args{
				&model.Friend{
					FromUser: uuid.UUID{},
					ToUser:   fakeToUUID,
				},
			},
			wantErr: true,
		},

		{
			name: "Valid testcase",
			args: args{
				&model.Friend{
					FromUser: fakeFromUUID,
					ToUser:   uuid.UUID{},
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
			if err := mw.Save(tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("validationMiddleware.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validationMiddleware_MakeFriend(t *testing.T) {
	serviceMock := &ServiceMock{
		MakeFriendFunc: func(from, to uuid.UUID) error {
			return nil
		},
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type args struct {
		from uuid.UUID
		to   uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid testcase",
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},

		{
			name: "FromUser UUID empty",
			args: args{
				from: uuid.UUID{},
				to:   fakeToUUID,
			},
			wantErr: true,
		},

		{
			name: "Valid testcase",
			args: args{
				from: fakeFromUUID,
				to:   uuid.UUID{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}
			if err := mw.MakeFriend(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("validationMiddleware.MakeFriend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validationMiddleware_UnFriend(t *testing.T) {
	serviceMock := &ServiceMock{
		UnFriendFunc: func(from, to uuid.UUID) error {
			return nil
		},
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type args struct {
		from uuid.UUID
		to   uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid testcase",
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},

		{
			name: "FromUser UUID empty",
			args: args{
				from: uuid.UUID{},
				to:   fakeToUUID,
			},
			wantErr: true,
		},

		{
			name: "Valid testcase",
			args: args{
				from: fakeFromUUID,
				to:   uuid.UUID{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}
			if err := mw.UnFriend(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("validationMiddleware.UnFriend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validationMiddleware_Accept(t *testing.T) {
	serviceMock := &ServiceMock{
		AcceptFunc: func(from, to uuid.UUID) error {
			return nil
		},
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type args struct {
		from uuid.UUID
		to   uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid testcase",
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},

		{
			name: "FromUser UUID empty",
			args: args{
				from: uuid.UUID{},
				to:   fakeToUUID,
			},
			wantErr: true,
		},

		{
			name: "Valid testcase",
			args: args{
				from: fakeFromUUID,
				to:   uuid.UUID{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}
			if err := mw.Accept(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("validationMiddleware.Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validationMiddleware_Reject(t *testing.T) {
	serviceMock := &ServiceMock{
		RejectFunc: func(from, to uuid.UUID) error {
			return nil
		},
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type args struct {
		from uuid.UUID
		to   uuid.UUID
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid testcase",
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},

		{
			name: "FromUser UUID empty",
			args: args{
				from: uuid.UUID{},
				to:   fakeToUUID,
			},
			wantErr: true,
		},

		{
			name: "Valid testcase",
			args: args{
				from: fakeFromUUID,
				to:   uuid.UUID{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}
			if err := mw.Reject(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("validationMiddleware.Reject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validationMiddleware_GetFriends(t *testing.T) {
	serviceMock := &ServiceMock{
		GetFriendsFunc: func(userID uuid.UUID) ([]model.Friend, error) {
			return nil, nil
		},
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantRes []model.Friend
		wantErr bool
	}{
		{
			name: "valid testcase 1",
			args: args{
				userID: fakeFromUUID,
			},
			wantRes: nil,
		},
		{
			name: "valid testcase 2",
			args: args{
				userID: fakeToUUID,
			},
			wantRes: nil,
		},
		{
			name: "empty uuid",
			args: args{
				userID: uuid.UUID{},
			},
			wantRes: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := validationMiddleware{
				Service: serviceMock,
			}
			gotRes, err := mw.GetFriends(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validationMiddleware.GetFriends() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("validationMiddleware.GetFriends() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
