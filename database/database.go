package database

import (
	"log"

	"github.com/jinzhu/gorm"
)

//DB is database
type DB struct {
	*gorm.DB
	url string
}

//New database
func New(url string) *DB {
	db := DB{url: url}
	db.initialize()
	db.migrate()
	return &db
}

func (db *DB) initialize() {
	if postgresDB, err := gorm.Open("postgres", db.url); err != nil {
		log.Fatal(err)
	} else {
		db.DB = postgresDB
	}
}

func (db *DB) migrate() {
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&FileHistory{})
}
