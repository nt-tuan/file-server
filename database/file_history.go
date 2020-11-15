package database

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

// Errors
var (
	ErrCreateHistory = errors.New("create-file-history-error")
)

//AddFileHistory to db
func (db *DB) AddFileHistory(file *File, action FileAction, fullname string, backupFullname *string) error {
	fileHistory := NewFileHistory(file, action, fullname, backupFullname)
	if err := db.Model(&FileHistory{}).
		Create(&fileHistory).
		Error; err != nil {
		return fmt.Errorf("%v: %w", err.Error(), ErrCreateHistory)
	}
	return nil
}

// GetFileHistoryRecords return changes of a file
func (db *DB) GetFileHistoryRecords(fileID uint) ([]FileHistory, error) {
	var files []FileHistory
	if err := db.Model(&FileHistory{}).Find(&files, &FileHistory{FileID: fileID}).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (db *DB) GetDeletedFileByID(id uint) (*FileHistory, error) {
	var file FileHistory
	if err := db.First(&file, FileHistory{Model: gorm.Model{ID: id}}).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// GetDeletedFiles return FileHistory which has DeleteAction
func (db *DB) GetDeletedFiles() ([]FileHistory, error) {
	var files []FileHistory
	if err := db.Model(&FileHistory{}).
		Where(FileHistory{ActionType: DeleteAction}).
		Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// RestoreDeletedFile will update DeleteAction to RestoreAction
func (db *DB) RestoreDeletedFile(fileHistory FileHistory) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Model(&FileHistory{Model: gorm.Model{ID: fileHistory.ID}}).
			Update("BackupFullname", nil).Error; err != nil {
			// return any error will rollback
			return err
		}
		if err := tx.Create(&FileHistory{
			Fullname:   fileHistory.Fullname,
			ActionType: RestoreAction,
			FileID:     fileHistory.FileID,
		}).Error; err != nil {
			return err
		}
		// return nil will commit the whole transaction
		return nil
	})
}
