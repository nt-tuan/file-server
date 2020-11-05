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
	local.ValidExts = []string{PngExt, SvgExt}
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
	serverPath, clientPath, err := lc.correctFileName(fileName)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(serverPath), os.ModePerm); err != nil {
		return "", err
	}
	out, err := os.Create(serverPath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, reader)
	return clientPath, err
}

// AddFile from fileheader
func (lc *Storage) AddFile(reader io.Reader, fileName string) (*database.File, error) {
	clientPath, err := lc.physicalAddFile(reader, fileName)
	if err != nil {
		return nil, err
	}
	// Save new file to database if this file created successfully
	fileModel := database.File{Fullname: clientPath}
	err = lc.db.CreateFile(&fileModel)

	// If failed to save to database, delete the file
	if err != nil {
		lc.DeleteFile(&fileModel)
		return nil, err
	}
	return &fileModel, nil
}

// ReplaceFile in storage
func (lc *Storage) ReplaceFile(dbFile *database.File, reader io.Reader) (string, error) {
	// Delete physical file
	log.Printf("Try delete file %s", dbFile.Fullname)
	backupPath, err := lc.physicalDeleteFile(dbFile)
	if err != nil {
		return "", err
	}
	// Create new physical file
	log.Printf("Try add file %s", dbFile.Fullname)
	_, err = lc.physicalAddFile(reader, dbFile.Fullname)
	if err != nil {
		return "", err
	}
	return backupPath, nil
}

//RenameFile in storage
func (lc *Storage) RenameFile(clientPath, newName string) (string, error) {
	// Find the file in database
	file, err := lc.db.GetFileByName(clientPath)
	if err != nil {
		return "", err
	}

	oldPsPath := lc.GetPhysicalWorkingPath(clientPath)
	newPsPath := lc.GetPhysicalWorkingPath(newName)

	if fileExists(newPsPath) {
		return "", ErrFileExisted
	}
	if err := os.MkdirAll(filepath.Dir(newPsPath), os.ModePerm); err != nil {
		return "", err
	}
	err = os.Rename(oldPsPath, newPsPath)
	if err != nil {
		return "", err
	}

	// Save rename action to database.
	// If failed to save action, rename to the origin one
	if err := lc.db.RenameFile(file, newName); err != nil {
		os.Rename(newPsPath, oldPsPath)
		return "", err
	}
	return newName, nil
}

// RollbackRenameFile will try to rollback of action renamefile
func (lc *Storage) RollbackRenameFile(path, newName string) (err error) {
	var dbFile *database.File
	dbFile, err = lc.db.GetFileByName(newName)
	if err != nil {
		return
	}

	err = lc.db.RenameFile(dbFile, path)
	return
}

func (lc *Storage) physicalDeleteFile(file *database.File) (string, error) {
	clientFileDir := filepath.Dir(file.Fullname)
	// Copy the file to history zone
	psPath := lc.GetPhysicalWorkingPath(file.Fullname)
	hsPath := lc.GetPhysicalHistoricalPath(file.Fullname)
	dst, err := copyFile(psPath, hsPath, true)
	if err != nil {
		return "", err
	}

	// Delete the actual file in working zone
	err = os.Remove(psPath)
	if err != nil {
		// If can not save the file, copy the from the history zone to working zone and remove the file in history zone
		if _, cfErr := copyFile(hsPath, psPath, false); cfErr == nil {
			os.Remove(dst)
		}
		// Delete copyfile when removing the file is getting error
		os.Remove(dst)
		return "", err
	}
	dstBase := filepath.Base(dst)
	backupPath := filepath.Join(clientFileDir, dstBase)
	return backupPath, nil
}

// DeleteFile will copy the file to history zone, then remove the file in working zone
// return the backup file and error if exists
func (lc *Storage) DeleteFile(file *database.File) error {
	backupPath, err := lc.physicalDeleteFile(file)
	if err != nil {
		return err
	}
	return lc.db.DeleteFile(file, backupPath)
}

// RollbackDeleteFile will try to rollback action deletefile
// Get the backup file and copy that file to working zone
func (lc *Storage) RollbackDeleteFile(fileName, bkFile string) (err error) {
	bkPath := filepath.Join(lc.HistoryDir, filepath.Base(bkFile))
	path := filepath.Join(filepath.Base(fileName), lc.WorkingDir)
	if _, err = copyFile(bkPath, path, false); err != nil {
		return
	}
	dbFile := database.File{Fullname: fileName}
	err = lc.db.CreateFile(&dbFile)
	return
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

		return lc.db.CreateFile(&database.File{Fullname: localPath})
	})
}

// GetFilePath from filename
func (lc *Storage) GetFilePath(filename string) string {
	return filepath.Join(lc.WorkingDir, filename)
}
