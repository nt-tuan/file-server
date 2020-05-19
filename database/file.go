package database

import (
	"path/filepath"

	"github.com/jinzhu/gorm"
)

// region gets

// GetFiles has all of specified tags
func (db *DB) GetFiles(tags []string, page, size uint, orders []string) ([]File, error) {
	var files []File
	tempDB := db.Model(&File{}).
		Preload("Tags")
	if tags != nil && len(tags) > 0 {
		tempDB = tempDB.
			Joins("JOIN file_tags ON file_tags.file_id = files.id").
			Where("file_tags.tag_id in ?", tags)
	}

	if orders != nil {
		for _, od := range orders {
			tempDB = tempDB.Order(od)
		}
	}

	if err := tempDB.
		Offset(size * page).
		Limit(size).
		Find(&files).
		Error; err != nil {
		return nil, err
	}
	return files, nil
}

//GetFileByName return File
func (db *DB) GetFileByName(filename string) (file *File, err error) {
	file = &File{Fullname: filename}
	err = db.Model(&File{}).
		First(file).Error
	return
}

//GetFileByID return file
func (db *DB) GetFileByID(id uint) (file *File, err error) {
	file = &File{Model: gorm.Model{ID: id}}
	err = db.Model(&File{}).
		First(file).Error
	return
}

//CountFiles has all of specified tags
func (db *DB) CountFiles(tags []string) (uint, error) {
	var count uint
	if err := db.Model(&File{}).
		Preload("Tags").
		Joins("JOIN file_tags ON file_tags.file_id = files.id").
		Where("file_tags.tag_id in ?", tags).
		Count(&count).
		Error; err != nil {
		return 0, err
	}
	return count, nil
}

// endregion

// CreateFile to database
func (db *DB) CreateFile(file *File) error {
	file.ExtractParts()
	if err := db.Model(&File{}).
		Create(file).
		Error; err != nil {
		return err
	}
	return db.AddFileHistory(file, CreateAction, file.Fullname)
}

//RenameFile in database
func (db *DB) RenameFile(file *File, newName string) error {
	file.Fullname = newName
	file.ExtractParts()
	if err := db.Model(&File{}).
		Updates(file).
		Error; err != nil {
		return err
	}
	return db.AddFileHistory(file, RenameAction, file.Fullname)
}

//DeleteFile in database
func (db *DB) DeleteFile(file *File, backup string) error {
	if err := db.Model(&File{}).
		Delete(file).Error; err != nil {
		return err
	}
	return db.AddFileHistory(file, DeleteAction, backup)
}

// UntrackDeleteFile will try to delete file in db without leaving a bread piece
func (db *DB) UntrackDeleteFile(path string) {
	name := filepath.Base(path)
	db.Model(&File{}).Delete(&File{Fullname: name})
}
