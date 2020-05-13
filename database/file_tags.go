package database

//AddTag to a file
func (db *DB) AddTag(file *File, tag string) error {
	if err := db.
		Model(&File{}).
		Association("Tags").
		Append(&Tag{ID: tag}).
		Error; err != nil {
		return err
	}
	return nil
}

//RemoveTag from a file
func (db *DB) RemoveTag(file *File, tag string) error {
	if err := db.Model(&File{}).
		Association("Tags").
		Delete(&Tag{ID: tag}).
		Error; err != nil {
		return err
	}
	return nil
}
