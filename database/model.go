package database

import (
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
)

// FileActions
var (
	CreateAction = "Created"
	RenameAction = "Renamed"
	DeleteAction = "Deleted"
)

// File table
type File struct {
	gorm.Model
	Fullname      string `gorm:"unique_index"`
	NamePart      string
	ExtensionPart *string
	Tags          []Tag `gorm:"many2many:file_tags;"`
	FileHistories []FileHistory
}

// FileHistory table store history of file changing
type FileHistory struct {
	gorm.Model
	Fullname      string
	NamePart      string
	ExtensionPart *string
	ActionType    string
	FileID        uint
	File          *File
}

//NewFileHistory created from File and action
func NewFileHistory(f *File, action string, backup string) *FileHistory {
	var h = FileHistory{}
	h.Fullname = backup
	h.ExtractParts()
	h.ActionType = action
	h.FileID = f.ID
	return &h
}

// ExtractParts from file
func (f *FileHistory) ExtractParts() {
	ext := filepath.Ext(f.Fullname)
	f.ExtensionPart = &ext
	parts := strings.SplitN(filepath.Base(f.Fullname), ".", 2)
	if len(parts) < 1 {
		return
	}
	f.NamePart = parts[0]
}

// ExtractParts from file
func (f *File) ExtractParts() {
	ext := filepath.Ext(f.Fullname)
	f.ExtensionPart = &ext
	parts := strings.SplitN(filepath.Base(f.Fullname), ".", 2)
	if len(parts) < 1 {
		return
	}
	f.NamePart = parts[0]
}

// Tag table
type Tag struct {
	ID    string `gorm:"primary_key"`
	Files []File `gorm:"many2many:file_tags;"`
}
