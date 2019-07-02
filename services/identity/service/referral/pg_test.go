package referral

import (
	"reflect"
	"testing"

	"github.com/dwarvesf/yggdrasil/services/identity/model"
	"github.com/dwarvesf/yggdrasil/services/identity/util/testutil"
	"github.com/jinzhu/gorm"

	uuid "github.com/satori/go.uuid"
)

func Test_pgService_Save(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		panic(err)
	}

	fakeUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		o *model.Referral
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Testcase 1",
			fields: fields{
				db,
			},
			args: args{
				&model.Referral{
					FromUserID:  fakeUUID,
					ToUserEmail: "meocon@meocon.com",
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

func Test_pgService_DeleteReferralWithCode(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		panic(err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		panic(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid code",
			fields: fields{
				db,
			},
			args: args{
				code: "123456",
			},
		},
		{
			name: "Invalid code",
			fields: fields{
				db,
			},
			args: args{
				code: "123406",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.DeleteReferralWithCode(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("pgService.DeleteReferralWithCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgService_Get(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		panic(err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		panic(err)
	}

	wantRes := model.Referral{}
	if err := db.First(&wantRes); err.Error != nil {
		panic(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		q *Query
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Referral
		wantErr bool
	}{
		{
			name: "code valid",
			fields: fields{
				db,
			},
			args: args{
				&Query{
					Code: "123456",
				},
			},
			want: wantRes,
		},
		{
			name: "code invalid",
			fields: fields{
				db,
			},
			args: args{
				&Query{
					Code: "123450",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			got, err := s.Get(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
