package db

import (
	"encoding/json"

	consul "github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/dwarvesf/yggdrasil/toolkit"
)

type PG struct {
	User     string
	Password string
	DB       string
}

func New(c *consul.Client) (*gorm.DB, func()) {
	pgAddress, pgPort, err := toolkit.GetServiceAddress(c, "postgres")
	if err != nil {
		panic(err)
	}

	v, err := toolkit.GetConsulValueFromKey(c, "db-scheduler")
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
