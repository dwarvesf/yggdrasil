package follow

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/networks/model"
	"github.com/dwarvesf/yggdrasil/networks/util/testutil"
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
		o *model.Follow
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid test case",
			fields: fields{
				db,
			},
			args: args{
				&model.Follow{
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

func Test_pgService_UnFollow(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		t.Errorf("Migarate table follow error = %v", err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		t.Errorf("Create sample data error = %v", err)
	}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")
	errorUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceeee")

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		fromUser uuid.UUID
		toUser   uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid test case 1",
			fields: fields{
				db,
			},
			args: args{
				fromUser: fakeFromUUID,
				toUser:   fakeToUUID,
			},
		},
		{
			name: "FromUser invalid",
			fields: fields{
				db,
			},
			args: args{
				fromUser: errorUUID,
				toUser:   fakeFromUUID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.UnFollow(tt.args.fromUser, tt.args.toUser); (err != nil) != tt.wantErr {
				t.Errorf("pgService.UnFollow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgService_FindAll(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		t.Errorf("Migarate table follow error = %v", err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		t.Errorf("Create sample data error = %v", err)
	}

	follow := []model.Follow{}
	if err := db.Find(&follow); err.Error != nil {
		t.Errorf("Get sample follow error = %v", err)
	}

	emptyFollow := []model.Follow{}

	fakeFromUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebf")
	fakeToUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceebc")
	errorUUID, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7ceeee")

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
		want    []model.Follow
		wantErr bool
	}{
		{
			name: "Valid test case from user",
			fields: fields{
				db,
			},
			args: args{
				&Query{
					FromUser: fakeFromUUID,
				},
			},
			want: follow,
		},
		{
			name: "Valid test case to user",
			fields: fields{
				db,
			},
			args: args{
				&Query{
					ToUser: fakeToUUID,
				},
			},
			want: follow,
		},
		{
			name: "Valid test case to user and from user",
			fields: fields{
				db,
			},
			args: args{
				&Query{
					ToUser:   fakeToUUID,
					FromUser: fakeFromUUID,
				},
			},
			want: follow,
		},
		{
			name: "To User not exists",
			fields: fields{
				db,
			},
			args: args{
				&Query{
					ToUser: errorUUID,
				},
			},
			want: emptyFollow,
		},
		{
			name: "From User not exists",
			fields: fields{
				db,
			},
			args: args{
				&Query{
					FromUser: errorUUID,
				},
			},
			want: emptyFollow,
		},
		{
			name: "To user and From user not exists",
			fields: fields{
				db,
			},
			args: args{
				&Query{
					ToUser:   errorUUID,
					FromUser: errorUUID,
				},
			},
			want: emptyFollow,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			got, err := s.FindAll(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgService.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
