package user

import (
	"errors"
	"testing"

	"github.com/dwarvesf/yggdrasil/identity/model"
)

func Test_validationMiddleware_Save(t *testing.T) {
	serviceMock := &ServiceMock{
		SaveFunc: func(user *model.User) error {
			return nil
		},
		GetFunc: func(userQuery *UserQuery) (*model.User, error) {
			return nil, errors.New("meo")
		},
	}

	type args struct {
		r *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Username valid",
			args: args{
				&model.User{
					LoginType: "username",
					Username:  "meocon",
					Password:  "123",
				},
			},
		},
		{
			name: "Username empty",
			args: args{
				&model.User{
					LoginType: "username",
					Username:  "",
					Password:  "123",
				},
			},
			wantErr: true,
		},
		{
			name: "Email invalid",
			args: args{
				&model.User{
					LoginType: "email",
					Email:     "meocon",
					Password:  "123",
				},
			},
			wantErr: true,
		},
		{
			name: "Email valid",
			args: args{
				&model.User{
					LoginType: "email",
					Email:     "meo@meowmeow.com",
					Password:  "123",
				},
			},
		},
		{
			name: "Empty email",
			args: args{
				&model.User{
					LoginType: "email",
					Email:     "",
					Password:  "123",
				},
			},
			wantErr: true,
		},
		{
			name: "PhoneNumber invalid",
			args: args{
				&model.User{
					LoginType: "email",
					Email:     "Minh bat chuoc loai meo kieu nha",
					Password:  "123",
				},
			},
			wantErr: true,
		},
		{
			name: "Phonenumber valid",
			args: args{
				&model.User{
					LoginType:   "phone_number",
					PhoneNumber: "0637893034",
					Password:    "123",
				},
			},
		},
		{
			name: "Empty phonenumber",
			args: args{
				&model.User{
					LoginType:   "phone_number",
					PhoneNumber: "",
					Password:    "123",
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
				t.Errorf("validationMiddleware.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
