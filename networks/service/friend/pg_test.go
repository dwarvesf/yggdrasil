package friend

import (
	"reflect"
	"testing"
	"time"

	"github.com/dwarvesf/yggdrasil/networks/model"
	"github.com/dwarvesf/yggdrasil/networks/util/testutil"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func Test_pgService_Save(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		t.Errorf("Migarate table follow error = %v", err)
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		o *model.Friend
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid",
			fields: fields{
				db,
			},
			args: args{
				&model.Friend{
					FromUser: fakeFromUUID,
					ToUser:   fakeToUUID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.Save(tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("pgService.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgService_MakeFriend(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		t.Errorf("Migarate table follow error = %v", err)
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		from uuid.UUID
		to   uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "New friend",
			fields: fields{
				db,
			},
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},
		{
			name: "Not accepted friend",
			fields: fields{
				db,
			},
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.MakeFriend(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("pgService.MakeFriend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgService_UnFriend(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		t.Errorf("Migarate table data error = %v", err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		t.Errorf("Create sample data error = %v", err)
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")
	missUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceaaa")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		from uuid.UUID
		to   uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid testcase",
			fields: fields{
				db,
			},
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},
		{
			name: "To UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: fakeFromUUID,
				to:   missUUID,
			},
			wantErr: true,
		},
		{
			name: "From UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: missUUID,
				to:   fakeToUUID,
			},
			wantErr: true,
		},
		{
			name: "UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: missUUID,
				to:   missUUID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.UnFriend(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("pgService.UnFriend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgService_Accept(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		t.Errorf("Migarate table data error = %v", err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		t.Errorf("Create sample data error = %v", err)
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")
	missUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceaaa")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		from uuid.UUID
		to   uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid testcase",
			fields: fields{
				db,
			},
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},
		{
			name: "To UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: fakeFromUUID,
				to:   missUUID,
			},
			wantErr: true,
		},
		{
			name: "From UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: missUUID,
				to:   fakeToUUID,
			},
			wantErr: true,
		},
		{
			name: "UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: missUUID,
				to:   missUUID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.Accept(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("pgService.Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgService_Reject(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		t.Errorf("Migarate table data error = %v", err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		t.Errorf("Create sample data error = %v", err)
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")
	missUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceaaa")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		from uuid.UUID
		to   uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid testcase",
			fields: fields{
				db,
			},
			args: args{
				from: fakeFromUUID,
				to:   fakeToUUID,
			},
		},
		{
			name: "To UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: fakeFromUUID,
				to:   missUUID,
			},
			wantErr: true,
		},
		{
			name: "From UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: missUUID,
				to:   fakeToUUID,
			},
			wantErr: true,
		},
		{
			name: "UUID has not existed yet",
			fields: fields{
				db,
			},
			args: args{
				from: missUUID,
				to:   missUUID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.Reject(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("pgService.Reject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgService_GetFriends(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		t.Errorf("Migarate table data error = %v", err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		t.Errorf("Create sample data error = %v", err)
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")
	missUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceaaa")

	res := model.Friend{}
	if err := db.Find(&res); err.Error != nil {
		t.Errorf("Accept err = %v", err)
	}

	// make friend accept request
	if err := db.Model(&res).Update(map[string]interface{}{"accepted_at": time.Now(),
		"rejected_at":  time.Time{},
		"requested_at": time.Time{}}).Error; err != nil {
		t.Errorf("Accept err = %v", err)
	}

	wantData := []model.Friend{}
	if err := db.Find(&wantData); err.Error != nil {
		t.Errorf("Get want data err = %v", err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Friend
		wantErr bool
	}{
		{
			name: "FromUUID exist",
			fields: fields{
				db,
			},
			args: args{
				userID: fakeFromUUID,
			},
			want: wantData,
		},
		{
			name: "ToUUID exist",
			fields: fields{
				db,
			},
			args: args{
				userID: fakeToUUID,
			},
			want: wantData,
		},
		{
			name: "UUID not exist",
			fields: fields{
				db,
			},
			args: args{
				userID: missUUID,
			},
			want: []model.Friend{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			got, err := s.GetFriends(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgService.GetFriend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgService.GetFriend() = %v, want %v", got, tt.want)
			}
		})
	}
}
