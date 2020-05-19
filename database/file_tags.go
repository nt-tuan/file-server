package database

//AddTag to a file
func (db *DB) AddTag(file *File, tagstr string) error {
	tag := Tag{ID: tagstr}
	db.Create(&tag)
	if err := db.
		Model(file).
		Association("Tags").
		Append(&tag).
		Error; err != nil {
		return err
	}
	return nil
}

//RemoveTag from a file
func (db *DB) RemoveTag(file *File, tag string) error {
	if err := db.Model(file).
		Association("Tags").
		Delete(&Tag{ID: tag}).
		Error; err != nil {
		return err
	}
	return nil
}
