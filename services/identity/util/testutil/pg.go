package testutil

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/satori/go.uuid"

	"github.com/dwarvesf/yggdrasil/services/identity/model"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID           string
	Username     string
	Phong_Number string
	Email        string
	Login_Type   string
	Password     string
	Status       string
}

// DB info
const (
	dbHost     = "localhost"
	dbPort     = 5439
	dbUser     = "postgres"
	dbPassword = "123"
	dbName     = "test"
)

func MigrateTable(db *gorm.DB) error {
	db.CreateTable(model.User{})
	db.CreateTable(model.Referral{})
	return nil
}

func CreateSampleData(db *gorm.DB) error {
	uuid, _ := uuid.FromString("42fd427e-0753-4a00-aed5-483a0c7c62ab")
	user := model.User{
		ID:        uuid,
		LoginType: "username",
		Username:  "meocon",
		Password:  "123",
		Status:    2,
	}

	referral := model.Referral{
		ID:          uuid,
		Code:        "123456",
		FromUserID:  uuid,
		ToUserEmail: "meocon@meocon.org",
	}

	db.Create(&user)
	db.Create(&referral)
	return nil
}

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
