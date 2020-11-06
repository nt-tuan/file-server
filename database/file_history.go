package database

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrCreateHistory = errors.New("create-file-history-error")
)

//AddFileHistory to db
func (db *DB) AddFileHistory(file *File, action string, fullname string, backupFullname *string) error {
	fileHistory := NewFileHistory(file, action, fullname, backupFullname)
	if err := db.Model(&FileHistory{}).
		Create(&fileHistory).
		Error; err != nil {
		return fmt.Errorf("%v: %w", err.Error(), ErrCreateHistory)
	}
	return nil
}
