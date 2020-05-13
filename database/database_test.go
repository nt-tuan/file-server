package database

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *DB

func setup() {
	db = &DB{url: "postgres://file-server:@:5432/file-server?sslmode=disable"}
	db.initialize()
	db.teardown()
	db.migrate()
	db.LogMode(true)
}

func (db *DB) teardown() {
	db.DropTableIfExists(&FileHistory{})
	db.DropTableIfExists(&File{})
	db.DropTableIfExists(&Tag{})
}

func TestAddFile(t *testing.T) {
	newFile := &File{Fullname: "test.txt"}
	newFile2 := &File{Fullname: "test.txt"}
	if err := db.AddFile(newFile); err != nil {
		t.Error(err)
	}
	if err := db.AddFile(newFile2); err == nil {
		t.Error(err)
	}
	if err := db.RenameFile(newFile, "hhee"); err != nil {
		t.Error(err)
	}
	if err := db.ReplaceFile(newFile, "/removed/hehehe"); err != nil {
		t.Error(err)
	}
	if err := db.DeleteFile(&File{Model: gorm.Model{ID: newFile.ID}}, "/deleted/zzz.txt"); err != nil {
		t.Error(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}
