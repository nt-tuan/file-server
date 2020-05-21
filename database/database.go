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
	db.Initialize()
	db.Migrate()
	return &db
}

//NewClean database
func NewClean(url string) *DB {
	db := &DB{url: url}
	db.Initialize()
	db.Teardown()
	db.Migrate()
	return db
}

//SetURL of database
func (db *DB) SetURL(url string) {
	db.url = url
}

//Initialize database
func (db *DB) Initialize() {
	if postgresDB, err := gorm.Open("postgres", db.url); err != nil {
		log.Fatal(err)
	} else {
		db.DB = postgresDB
	}
}

//Teardown database
func (db *DB) Teardown() {
	db.DropTableIfExists(&FileHistory{})
	db.DropTableIfExists("file_tags")
	db.DropTableIfExists("FileTag")
	db.DropTableIfExists(&File{})
	db.DropTableIfExists(&Tag{})
}
