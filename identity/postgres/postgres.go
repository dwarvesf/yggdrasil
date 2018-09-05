package postgres

import (
	"github.com/jinzhu/gorm"
)

func New(ds string) (*gorm.DB, func() error) {
	db, err := gorm.Open("postgres", ds)
	if err != nil {
		panic(err)
	}
	return db, db.Close
}

// NewFake returns an empty gorm.DB struct for the sake of demonstration,
// you should use New() and remove this function afterward.
func NewFake(ds string) (*gorm.DB, func() error) {
	return &gorm.DB{}, func() error { return nil }
}
