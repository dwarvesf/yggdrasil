package referral

import (
	"testing"

	"github.com/dwarvesf/yggdrasil/services/identity/model"
	uuid "github.com/satori/go.uuid"
)

func Test_validationMiddleware_DeleteReferralWithCode(t *testing.T) {
	serviceMock := &ServiceMock{
		DeleteReferralWithCodeFunc: func(code string) error {
			return nil
		},
	}

	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid test case",
			args: args{
				code: "123456",
			},
		},
		{
			name: "empty code",
			args: args{
				code: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := &validationMiddleware{
				Service: serviceMock,
			}
			if err := mw.DeleteReferralWithCode(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("validationMiddleware.DeleteReferralWithCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validationMiddleware_Save(t *testing.T) {
	serviceMock := &ServiceMock{
		SaveFunc: func(o *model.Referral) error {
			return nil
		},
	}

	fakeUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")

	type args struct {
		o *model.Referral
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid test case",
			args: args{
				&model.Referral{
					FromUserID:  fakeUUID,
					ToUserEmail: "meocon@meowmeow.com",
				},
			},
		},
		{
			name: "Email empty",
			args: args{
				&model.Referral{
					FromUserID: fakeUUID,
				},
			},
			wantErr: true,
		},
		{
			name: "uuid empty",
			args: args{
				&model.Referral{
					FromUserID: fakeUUID,
				},
			},
			wantErr: true,
		},
		{
			name: "email format invalid",
			args: args{
				&model.Referral{
					FromUserID:  fakeUUID,
					ToUserEmail: "meocon",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := &validationMiddleware{
				Service: serviceMock,
			}
			if err := mw.Save(tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("validationMiddleware.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
