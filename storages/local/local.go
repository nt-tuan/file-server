package localstorage

import (
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/thanhtuan260593/file-server/database"
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
	return &local
}

// AddFile from fileheader
func (lc *Storage) AddFile(reader io.Reader, fileName string) (string, error) {
	serverPath, clientPath, err := lc.correctFileName(fileName)
	if err != nil {
		return "", err
	}

	out, err := os.Create(serverPath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, reader)

	// Save new file to database if this file created successfully
	fileModel := database.File{Fullname: clientPath}
	err = lc.db.CreateFile(&fileModel)

	// If failed to save to database, delete the file
	if err != nil {
		lc.DeleteFile(clientPath)
		return "", err
	}
	return clientPath, nil
}

// RollbackNewFile will try to rollback of action newfile
func (lc *Storage) RollbackNewFile(path string) (err error) {
	lc.DeleteFile(path)
	return
}

// ReplaceFile in storage
func (lc *Storage) ReplaceFile(path string, file io.Reader) (string, error) {
	// Find file from database, if no file found, return error
	_, err := lc.db.GetFileByName(path)
	if err != nil {
		return "", err
	}

	// Delete the file
	var bkDelFile string
	bkDelFile, err = lc.DeleteFile(path)
	if err != nil {
		return "", err
	}

	// Create new file with the same name
	_, err = lc.AddFile(file, filepath.Base(path))

	// If failed to create file, rollback action delete file
	if err != nil {
		lc.RollbackDeleteFile(path, bkDelFile)
		return "", err
	}
	return bkDelFile, nil
}

//RenameFile in storage
func (lc *Storage) RenameFile(clientPath, newName string) (string, error) {
	// Find the file in database
	file, err := lc.db.GetFileByName(clientPath)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(clientPath)
	newPath := filepath.Join(dir, newName)

	oldPsPath := lc.GetPhysicalWorkingPath(clientPath)
	newPsPath := lc.GetPhysicalWorkingPath(newPath)

	err = os.Rename(oldPsPath, newPsPath)
	if err != nil {
		return "", err
	}

	// Save rename action to database.
	// If failed to save action, rename to the origin one
	if err := lc.db.RenameFile(file, newPath); err != nil {
		os.Rename(newPsPath, oldPsPath)
		return "", err
	}
	return newPath, nil
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

// DeleteFile will copy the file to history zone, then remove the file in working zone
// return the backup file and error if exists
func (lc *Storage) DeleteFile(fileName string) (string, error) {
	// Find the file from database, if no file found return error
	clientFileDir := filepath.Dir(fileName)

	file, err := lc.db.GetFileByName(fileName)
	if err != nil {
		return "", err
	}

	// Copy the file to history zone
	psPath := lc.GetPhysicalWorkingPath(fileName)
	hsPath := lc.GetPhysicalHistoricalPath(fileName)
	dst, err := copyFile(psPath, hsPath, true)
	if err != nil {
		return "", err
	}

	// Delete the actual file in working zone
	err = os.Remove(psPath)
	if err != nil {
		// Delete copyfile when removing the file is getting error
		os.Remove(dst)
		return "", err
	}

	// Remove the file in database and save its delete action
	dstBase := filepath.Base(dst)
	backupPath := filepath.Join(clientFileDir, dstBase)
	err = lc.db.DeleteFile(file, backupPath)

	// If can not save the file, copy the from the history zone to working zone and remove the file in history zone
	if err != nil {
		if _, cfErr := copyFile(hsPath, psPath, false); cfErr == nil {
			os.Remove(dst)
		}
		return "", err
	}
	return dst, nil
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

//GetResized2DImage in storage
func (lc *Storage) GetResized2DImage(filename string, width, height uint) (image.Image, error) {
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

	//Resize image
	resized := imaging.Resize(imageData, int(width), int(height), imaging.Lanczos)
	return resized, nil
}
