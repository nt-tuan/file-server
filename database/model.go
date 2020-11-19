package database

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
)

// FileAction is action type
type FileAction string

// FileActions
var (
	CreateAction  FileAction = "create"
	RenameAction  FileAction = "rename"
	ReplaceAction FileAction = "replace"
	DeleteAction  FileAction = "delete"
	RestoreAction FileAction = "retore"
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
	CreatedBy     string
	Width         int
	Height        int
	DiskSize      int64
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
	Fullname       string     `json:"fullname"`
	BackupFullname *string    `json:"backupFullname"`
	ActionType     FileAction `json:"actionType"`
	FileID         uint       `json:"fileID"`
	CreatedBy      string
	File           *File
}

// NewFile return a new file record
func NewFile(fullname string, width, height int, diskSize int64, user string) File {
	return File{Fullname: fullname, CreatedBy: user, Width: width, Height: height, DiskSize: diskSize}
}

//NewFileHistory created from File and action
func NewFileHistory(f *File, action FileAction, fullname string, backupFullname *string, user string) *FileHistory {
	var h = FileHistory{}
	h.Fullname = fullname
	h.BackupFullname = backupFullname
	h.ActionType = action
	h.FileID = f.ID
	h.CreatedBy = user
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
