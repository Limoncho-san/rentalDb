package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "host=localhost user=root password=root dbname=testingwithrentals sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}
