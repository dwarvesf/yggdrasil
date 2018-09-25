package user

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/yggdrasil/identity/model"
	"github.com/dwarvesf/yggdrasil/identity/util/testutil"
)

func Test_pgService_Save(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		panic(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		u *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Username valid",
			fields: fields{
				db,
			},
			args: args{
				&model.User{
					LoginType: "username",
					Username:  "meocon",
					Password:  "123",
					Status:    2,
				},
			},
		},
		{
			name: "Email valid",
			fields: fields{
				db,
			},
			args: args{
				&model.User{
					LoginType: "email",
					Email:     "meocon@xx.com",
					Password:  "123",
					Status:    2,
				},
			},
		},
		{
			name: "Phone valid",
			fields: fields{
				db,
			},
			args: args{
				&model.User{
					LoginType: "phone_number",
					Username:  "09090909099",
					Password:  "123",
					Status:    2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.Save(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("pgService.Save() error = %v, wantErr %v", err, tt.wantErr)
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

	user := model.User{}
	if err := db.First(&user); err.Error != nil {
		panic(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		userQuery *UserQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "Get with username valid",
			fields: fields{
				db,
			},
			args: args{
				&UserQuery{
					Username: "meocon",
				},
			},
			want: &user,
		},
		{
			name: "Get with username and logintype valid",
			fields: fields{
				db,
			},
			args: args{
				&UserQuery{
					Username:  "meocon",
					LoginType: "username",
				},
			},
			want: &user,
		},
		{
			name: "Get with username valid and login type invalid",
			fields: fields{
				db,
			},
			args: args{
				&UserQuery{
					Username:  "meocon",
					LoginType: "username",
				},
			},
			want: &user,
		},
		{
			name: "Get with username invalid",
			fields: fields{
				db,
			},
			args: args{
				&UserQuery{
					Username: "xxxx",
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
			got, err := s.Get(tt.args.userQuery)
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

func Test_pgService_MakeActive(t *testing.T) {
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
		user *model.User
	}

	user := model.User{}
	if err := db.First(&user); err.Error != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid test case",
			fields: fields{
				db,
			},
			args: args{
				user: &user,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			if err := s.MakeActive(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("pgService.MakeActive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pgService_Login(t *testing.T) {
	db, _, fi := testutil.CreateTestDatabase(t)
	defer fi()

	if err := testutil.MigrateTable(db); err != nil {
		panic(err)
	}

	if err := testutil.CreateSampleData(db); err != nil {
		panic(err)
	}

	user := model.User{}
	if err := db.First(&user); err.Error != nil {
		panic(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		loginType string
		identity  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "invalid test case",
			fields: fields{
				db,
			},
			args: args{
				loginType: "username",
				identity:  "meomeo",
			},
			wantErr: true,
		},
		{
			name: "valid test case",
			fields: fields{
				db,
			},
			args: args{
				loginType: "username",
				identity:  "meocon",
			},
			want: &user,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pgService{
				db: tt.fields.db,
			}
			got, err := s.Login(tt.args.loginType, tt.args.identity)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgService.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
