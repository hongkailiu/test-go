package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// OpenPostgres opens a postgres db
func OpenPostgres(postgres string) (db *gorm.DB, err error) {
	return gorm.Open("postgres", postgres)
}
