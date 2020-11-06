package database

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
)

// FileActions
var (
	CreateAction  = "Created"
	RenameAction  = "Renamed"
	ReplaceAction = "Replaced"
	DeleteAction  = "Deleted"
)

// Errors
var (
	ErrNotFound = errors.New("file-not-found")
)

// File table
type File struct {
	gorm.Model
	Fullname      string
	NamePart      string
	ExtensionPart *string
	Tags          []Tag `gorm:"many2many:file_tags;association_foreignkey:ID;foreignkey:ID"`
	FileHistories []FileHistory
}

// Tag table
type Tag struct {
	ID string `gorm:"primary_key:true"`
}

// FileHistory table store history of file changing
type FileHistory struct {
	gorm.Model
	Fullname       string
	BackupFullname *string
	ActionType     string
	FileID         uint
	File           *File
}

//NewFileHistory created from File and action
func NewFileHistory(f *File, action string, fullname string, backupFullname *string) *FileHistory {
	var h = FileHistory{}
	h.Fullname = fullname
	h.BackupFullname = backupFullname
	h.ActionType = action
	h.FileID = f.ID
	return &h
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
