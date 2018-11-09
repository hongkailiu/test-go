package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func OpenPostgres(postgres string) (db *gorm.DB, err error) {
	return gorm.Open("postgres", postgres)
}
