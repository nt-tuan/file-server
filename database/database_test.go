package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *DB
var dbURL = "postgres://image-server:@:54321/image-server?sslmode=disable"

func setup() {
	db = NewClean(dbURL)
}

func TestGetFiles(t *testing.T) {
	setup()
	files, err := db.GetFiles(make([]string, 0), 0, 10, make([]string, 0))
	if err != nil {
		t.Error(err)
		return
	}
	for _, f := range files {
		fmt.Println(f.Fullname)
	}
}

func TestAddFile(t *testing.T) {
	setup()
	newFile := &File{Fullname: "test.txt"}
	newFile2 := &File{Fullname: "test.txt"}
	if err := db.CreateFile(newFile); err != nil {
		t.Error(err)
	}
	if err := db.CreateFile(newFile2); err != nil {
		t.Error(err)
		return
	}
	if err := db.RenameFile(newFile, "hhee"); err != nil {
		t.Error(err)
		return
	}
	if err := db.DeleteFile(&File{Model: gorm.Model{ID: newFile.ID}}, "/deleted/zzz.txt"); err != nil {
		t.Error(err)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
