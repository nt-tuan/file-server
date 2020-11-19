package localstorage

import (
	"image"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ptcoffee/image-server/database"
)

// Storage file storage
type Storage struct {
	WorkingDir string
	HistoryDir string
	ValidExts  []string
	db         *database.DB
}

// NewStorage return new LocalStorage
func NewStorage(db *database.DB) *Storage {
	var local = Storage{}
	local.db = db
	local.ValidExts = []string{PNGExt, SVGExt, JPEGExt, JPGExt, GIFExt, WEBPExt}
	local.WorkingDir = DefaultWorkingDir
	local.HistoryDir = DefaultHistoryDir

	//Try get IMAGE_WORKING_DIR and IMAGE_HISTORY_DIR from os enviroment
	if w := os.Getenv("IMAGE_WORKING_DIR"); w != "" {
		local.WorkingDir = w
	}
	if w := os.Getenv("IMAGE_HISTORY_DIR"); w != "" {
		local.HistoryDir = w
	}

	if isInit := os.Getenv("INIT_SAMPLE_DATA"); isInit != "" {
		if v, err := strconv.ParseBool(isInit); err == nil && v {
			local.CreateMissingFiles()
		}
	}
	return &local
}

func (lc *Storage) physicalAddFile(reader io.Reader, fileName string) (string, error) {
	path, fullname, err := lc.correctFileName(fileName)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return "", err
	}
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, reader)
	return fullname, err
}

// AddFile will add file to storage
func (lc *Storage) AddFile(reader io.Reader, fullname string, user string) (*database.File, error) {
	if !lc.IsValidExt(filepath.Ext(fullname)) {
		return nil, ErrFileExtInvalid
	}
	clientPath, err := lc.physicalAddFile(reader, fullname)
	if err != nil {
		return nil, err
	}
	imageSize, _ := lc.GetImageSize(clientPath)
	diskSize, err := lc.GetFileSize(clientPath)
	if err != nil {
		return nil, err
	}
	// Save new file to database if this file created successfully
	file, err := lc.db.CreateFile(clientPath, imageSize.Width, imageSize.Height, diskSize, user)
	// If failed to save to database, delete the file
	if err != nil {
		lc.physicalDeleteFile(clientPath)
	}
	return file, err
}

// ReplaceFile will move the old file to trash and add a new file with the same name in storage
func (lc *Storage) ReplaceFile(file database.File, reader io.Reader, user string) (*string, error) {
	backupFullname, err := lc.physicalDeleteFile(file.Fullname)
	if err != nil {
		return nil, err
	}

	newName, err := lc.physicalAddFile(reader, file.Fullname)
	if err != nil {
		return nil, err
	}
	lc.db.ReplaceFile(&file, newName, backupFullname, user)
	return backupFullname, nil
}

//RenameFile in storage
func (lc *Storage) RenameFile(file database.File, newName string, user string) error {
	oldPath := lc.getPath(file.Fullname)
	newPath := lc.getPath(newName)

	if fileExists(newPath) {
		return ErrFileExisted
	}
	if err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm); err != nil {
		return err
	}
	err := moveFile(oldPath, newPath)
	if err != nil {
		return err
	}

	// Save rename action to database.
	// If failed to save action, rename to the origin one
	if err := lc.db.RenameFile(&file, newName, user); err != nil {
		if err := moveFile(newPath, oldPath); err != nil {
			log.Println(err)
		}
		return err
	}
	return nil
}

func (lc *Storage) physicalDeleteFile(fullname string) (*string, error) {
	path := lc.getPath(fullname)
	if fileExists(path) {
		backupFullname := lc.newBackupFullname(fullname)
		backupPath := lc.getBackupPath(backupFullname)
		err := moveFile(path, backupPath)
		if err != nil {
			return nil, err
		}
		return &backupFullname, nil
	}
	return nil, nil
}

// DeleteFile will copy the file to history zone, then remove the file in working zone
// return the backup file and error if exists
func (lc *Storage) DeleteFile(file database.File, user string) error {
	backupFullname, err := lc.physicalDeleteFile(file.Fullname)
	if err != nil {
		return err
	}
	return lc.db.DeleteFile(&file, backupFullname, user)
}

// RestoreDeletedFile will try to rollback action deletefile
// Get the backup file and copy that file to working zone
func (lc *Storage) RestoreDeletedFile(historyFile database.FileHistory, user string) (*database.File, error) {
	if historyFile.BackupFullname == nil {
		return nil, ErrFileNotFound
	}
	path := lc.getPath(historyFile.Fullname)
	if fileExists(path) {
		return nil, ErrFileExisted
	}
	backupPath := lc.getBackupPath(*historyFile.BackupFullname)
	if err := moveFile(backupPath, path); err != nil {
		return nil, err
	}
	imageSize, _ := lc.GetImageSize(historyFile.Fullname)
	diskSize, err := lc.GetFileSize(historyFile.Fullname)
	if err != nil {
		return nil, err
	}
	file, err := lc.db.CreateFile(historyFile.Fullname, imageSize.Width, imageSize.Height, diskSize, user)
	if err != nil {
		return nil, err
	}
	err = lc.db.RestoreDeletedFile(historyFile)
	if err != nil {
		return nil, err
	}
	return file, err
}

// GetImageBuffer return reader from filename
func (lc *Storage) GetImageBuffer(filename string) ([]byte, error) {
	path := lc.getPath(filename)
	return getImageBuffer(path)
}

// GetImage from filename
func (lc *Storage) GetImage(filename string) (image.Image, error) {
	var path = filepath.Join(lc.WorkingDir, filename)
	//Check if file extention is valid
	var ext = filepath.Ext(path)
	if !lc.IsValidExt(ext) {
		return nil, ErrFileExtInvalid
	}
	imageData, err := getImageFromPath(path)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

// CreateMissingFiles files
func (lc *Storage) CreateMissingFiles() {
	filepath.Walk(lc.WorkingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		localPath, err := filepath.Rel(lc.WorkingDir, path)
		if err != nil {
			return err
		}

		// skip if this file exists in database
		if _, err := lc.db.GetFileByName(localPath); err == nil {
			return nil
		}
		imageSize, _ := lc.GetImageSize(localPath)
		diskSize, err := lc.GetFileSize(localPath)
		if err != nil {
			return err
		}
		_, err = lc.db.CreateFile(localPath, imageSize.Width, imageSize.Height, diskSize, "")
		return err
	})
}
