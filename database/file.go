package database

import (
	"github.com/jinzhu/gorm"
)

func filterByTags(db *gorm.DB, tags []string) *gorm.DB {
	if tags != nil && len(tags) > 0 {
		return db.
			Joins("JOIN file_tags ON file_tags.file_id = files.id").
			Where("file_tags.tag_id in (?)", tags)
	}
	return db
}

// GetFiles has all of specified tags
func (db *DB) GetFiles(tags []string, offset *uint, limit *uint, orders []string) ([]File, error) {
	var files []File
	chain := db.Model(&File{}).
		Preload("Tags")
	chain = filterByTags(chain, tags)

	if orders != nil {
		for _, od := range orders {
			chain = chain.Order(od)
		}
	}
	if offset != nil {
		chain = chain.Offset(*offset)
	}
	if limit != nil {
		chain = chain.Limit(*limit)
	}
	if err := chain.
		Find(&files).
		Error; err != nil {
		return nil, err
	}
	return files, nil
}

// CountFiles return files number filtered by tags
func (db *DB) CountFiles(tags []string) (int, error) {
	chain := db.Model(&File{})
	chain = filterByTags(chain, tags)
	var count int
	if err := chain.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

//GetFileByName return File
func (db *DB) GetFileByName(filename string) (file *File, err error) {
	file = &File{Fullname: filename}
	err = db.Model(&File{}).
		Where(file).
		First(file).Error
	if file.ID == 0 {
		return nil, ErrNotFound
	}
	return
}

//GetFullFileByID return file
func (db *DB) GetFullFileByID(id uint) (file *File, err error) {
	file = &File{Model: gorm.Model{ID: id}}
	err = db.
		Preload("Tags").
		Model(&File{}).
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

// CreateFile to database
func (db *DB) CreateFile(fullname string, width, height int, diskSize int64, user string) (*File, error) {
	file := NewFile(fullname, width, height, diskSize, user)
	file.ExtractParts()
	if err := db.Model(&File{}).
		Create(&file).
		Error; err != nil {
		return nil, err
	}
	if err := db.AddFileHistory(&file, CreateAction, file.Fullname, nil, user); err != nil {
		return nil, err
	}
	return &file, nil
}

//RenameFile in database
func (db *DB) RenameFile(file *File, newName string, user string) error {
	file.Fullname = newName
	file.ExtractParts()
	if err := db.Model(&File{}).
		Updates(file).
		Error; err != nil {
		return err
	}
	return db.AddFileHistory(file, RenameAction, file.Fullname, nil, user)
}

// ReplaceFile in database
func (db *DB) ReplaceFile(file *File, newName string, backupName *string, user string) error {
	file.Fullname = newName
	file.ExtractParts()
	if err := db.Model(&File{}).
		Updates(file).Error; err != nil {
		return err
	}
	return db.AddFileHistory(file, ReplaceAction, newName, backupName, user)
}

//DeleteFile in database
func (db *DB) DeleteFile(file *File, backup *string, user string) error {
	if err := db.Model(&File{}).
		Delete(file).Error; err != nil {
		return err
	}
	return db.AddFileHistory(file, DeleteAction, file.Fullname, backup, user)
}
