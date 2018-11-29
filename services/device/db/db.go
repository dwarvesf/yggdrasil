package db

import (
	"encoding/json"

	consul "github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/yggdrasil/services/device/model"
	"github.com/dwarvesf/yggdrasil/toolkit"
)

//PG struct declare data to login postgres
type PG struct {
	User     string
	Password string
	DB       string
}

// Migrate use to migrate database
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Device{})
}

// New use to connect with database
func New(c *consul.Client) (*gorm.DB, func()) {
	pgAddress, pgPort, err := toolkit.GetServiceAddress(c, "postgres")
	if err != nil {
		panic(err)
	}

	v, err := toolkit.GetConsulValueFromKey(c, "db-identity")
	if err != nil {
		panic(err)
	}

	var pg PG
	err = json.Unmarshal([]byte(v), &pg)
	if err != nil {
		panic(err)
	}

	return toolkit.NewDB(toolkit.ToDS(pg.User, pg.Password, pgAddress, pgPort, "postgres", pg.DB))
}
