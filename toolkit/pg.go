package toolkit

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
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
