package device

import (
	"net/http"
	"testing"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/device/model"
	identityModel "github.com/dwarvesf/yggdrasil/services/identity/model"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

func Test_PGService_ValidateUser(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := toolkit.CreateTestDatabase(t)
	defer cleanup()

	err := toolkit.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	user := identityModel.User{}
	err = testDB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user by error %v", err)
	}

	fakeUserID, errParseUUID := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7c62ac")
	if errParseUUID != nil {
		t.Errorf("failed parse uuid from string")
	}

	type args struct {
		UserID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				UserID: user.ID,
			},
		},
		{
			name: "failed validate user by error invalid user id",
			args: args{
				UserID: fakeUserID,
			},
			wantErr: ErrInvalidUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			s := &pgService{
				db: testDB,
			}
			err := s.ValidateUser(tt.args.UserID)
			if err != nil && err != tt.wantErr {
				t.Errorf("pgService.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("pgService.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func Test_PGService_Create(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := toolkit.CreateTestDatabase(t)
	defer cleanup()

	err := toolkit.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	user := identityModel.User{}
	err = testDB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user by error %v", err)
	}

	device := model.Device{UserID: user.ID}
	err = testDB.Create(&device).Error
	if err != nil {
		t.Fatalf("Failed to create device by error %v", err)
	}

	type args struct {
		d *model.Device
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				&device,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: testDB,
			}
			err := s.Create(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgService.Create() error = %v", err)
			}
		})
	}

}

func Test_PGService_Get(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := toolkit.CreateTestDatabase(t)
	defer cleanup()

	err := toolkit.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	fakeDeviceID, errParseUUID := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7c62ab")
	if errParseUUID != nil {
		t.Errorf("failed parse uuid from string")
	}
	device := &model.Device{}
	err = testDB.Create(device).Error
	if err != nil {
		t.Fatalf("failed to create tabel test device by error %v", err)
	}

	type args struct {
		deviceID uuid.UUID
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		errorStatusCode int
	}{
		{
			name: "success",
			args: args{
				deviceID: device.ID,
			},
		},
		{
			name: "failed get device by error not found device id",
			args: args{
				deviceID: fakeDeviceID,
			},
			wantErr:         true,
			errorStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			s := &pgService{
				db: testDB,
			}
			_, err := s.Get(tt.args.deviceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgService.Get() error = %v", err)
			}
		})
	}
}

func Test_PGService_GetList(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := toolkit.CreateTestDatabase(t)
	defer cleanup()

	err := toolkit.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	user := identityModel.User{}
	err = testDB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user by error %v", err)
	}

	device := model.Device{UserID: user.ID}
	err = testDB.Create(&device).Error
	if err != nil {
		t.Fatalf("Failed to create device by error %v", err)
	}

	type args struct {
		Query Query
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				Query{
					UserID: device.UserID,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			s := &pgService{
				db: testDB,
			}

			_, err := s.GetList(tt.args.Query)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgService.GetList() error = %v", err)
			}
		})
	}
}

func Test_PGService_Update(t *testing.T) {
	t.Parallel()
	testDB, _, cleanup := toolkit.CreateTestDatabase(t)
	defer cleanup()

	err := toolkit.MigrateTables(testDB)
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	device := model.Device{}
	err = testDB.Create(&device).Error
	if err != nil {
		t.Fatalf("Failed to create device by error %v", err)
	}

	fakeDeviceID, errParseUUID := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7c62ac")
	if errParseUUID != nil {
		t.Errorf("failed parse uuid from string")
	}

	type args struct {
		Device *model.Device
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				&model.Device{
					ID:       device.ID,
					Status:   2,
					FCMToken: "token",
				},
			},
		},
		{
			name: "failed by not found id",
			args: args{
				&model.Device{
					ID:     fakeDeviceID,
					Status: 2,
				},
			},
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			s := &pgService{
				db: testDB,
			}
			err := s.Update(tt.args.Device)
			if err != nil && err != tt.wantErr {
				t.Errorf("pgService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("pgService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
