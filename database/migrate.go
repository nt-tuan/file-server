package database

//Migrate database
func (db *DB) Migrate() {
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&FileHistory{})
	db.Model(&FileHistory{}).AddForeignKey("file_id", "files(id)", "RESTRICT", "RESTRICT")
	db.Table("file_tags").AddForeignKey("file_id", "files(id)", "RESTRICT", "RESTRICT")
	db.Table("file_tags").AddForeignKey("tag_id", "tags(id)", "RESTRICT", "RESTRICT")
}
