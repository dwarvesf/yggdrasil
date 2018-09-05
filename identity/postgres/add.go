package postgres

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/yggdrasil/identity/service/add"
)

type addStore struct {
	db *gorm.DB
}

// NewAddStore ...
func NewAddStore(db *gorm.DB) *addStore {
	return &addStore{
		db: db,
	}
}

// Add just do a plus with 2 vars (X+Y) for the sake of demonstration, in reallity
// you would want to execute a DB query/transaction here
func (s *addStore) Add(ctx context.Context, arg *add.Add) (int, error) {
	return arg.X + arg.Y, nil
}
