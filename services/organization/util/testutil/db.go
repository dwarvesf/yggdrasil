package testutil

import (
	"fmt"
	"os"

	"github.com/dwarvesf/yggdrasil/services/organization/db"
	"github.com/jinzhu/gorm"
)

// DB info
const (
	dbHost     = "localhost"
	dbPort     = 5439
	dbUser     = "postgres"
	dbPassword = "123"
	dbName     = "test"
)

// GetDB for testing
func GetDB() *gorm.DB {
	testingHost := fmt.Sprintf("%s", dbHost)
	testingPort := fmt.Sprintf("%d", dbPort)

	if host := os.Getenv("POSTGRES_TESTING_HOST"); host != "" {
		testingHost = host
	}

	if port := os.Getenv("POSTGRES_TESTING_PORT"); port != "" {
		testingPort = port
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", testingHost, testingPort, dbUser, dbPassword, dbName)
	pgdb, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	clearDB(pgdb)
	db.Migrate(pgdb)

	return pgdb
}

func clearDB(pgdb *gorm.DB) {
	pgdb.Exec("DELETE FROM user_groups")
	pgdb.Exec("DELETE FROM user_organizations")
	pgdb.Exec("DELETE FROM groups")
	pgdb.Exec("DELETE FROM organizations")
}
