package toolkit

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jinzhu/gorm"

	deviceModel "github.com/dwarvesf/yggdrasil/services/device/model"
	identityModel "github.com/dwarvesf/yggdrasil/services/identity/model"
)

// NewDB return a db connection and a function to close DB
func NewDB(ds string) (*gorm.DB, func()) {
	db, err := gorm.Open("postgres", ds)
	if err != nil {
		panic(err)
	}

	return db, func() {
		err := db.Close()
		if err != nil {
			log.Println("Failed to close DB by error", err)
		}
	}
}

// ToDS return datasource of given values
func ToDS(user, password, address string, port int, dialect, name string) string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable", dialect, user, password, address, port, name)
}

const (
	dbHost     = "localhost"
	dbPort     = 5439
	dbUser     = "postgres"
	dbPassword = "123"
	dbName     = "test"
)

// CreateTestDatabase will create a test-database and test-schema
func CreateTestDatabase(t *testing.T) (*gorm.DB, string, func()) {
	testingHost := fmt.Sprintf("%s", dbHost)
	testingPort := fmt.Sprintf("%d", dbPort)
	if os.Getenv("POSTGRES_TESTING_HOST") != "" {
		testingHost = os.Getenv("POSTGRES_TESTING_HOST")
	}
	if os.Getenv("POSTGRES_TESTING_PORT") != "" {
		testingPort = os.Getenv("POSTGRES_TESTING_PORT")
	}
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", testingHost, testingPort, dbUser, dbPassword, dbName)
	db, dbErr := gorm.Open("postgres", connectionString)
	if dbErr != nil {
		t.Fatalf("Fail to create database. %s", dbErr.Error())
	}

	rand.Seed(time.Now().UnixNano())
	schemaName := "test" + strconv.FormatInt(rand.Int63(), 10)

	err := db.Exec("CREATE SCHEMA " + schemaName).Error
	if err != nil {
		t.Fatalf("Fail to create schema. %s", err.Error())
	}

	// set schema for current db connection
	err = db.Exec("SET search_path TO " + schemaName).Error
	if err != nil {
		t.Fatalf("Fail to set search_path to created schema. %s", err.Error())
	}

	return db, schemaName, func() {
		err := db.Exec("DROP SCHEMA " + schemaName + " CASCADE").Error
		if err != nil {
			t.Fatalf("Fail to drop database. %s", err.Error())
		}
	}
}

// MigrateTables migrate db with tables base by domain model
func MigrateTables(db *gorm.DB) error {
	return db.AutoMigrate(
		identityModel.User{},
		deviceModel.Device{},
	).Error
}
