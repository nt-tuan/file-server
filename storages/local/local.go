package local

import (
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/thanhtuan260593/file-server/database"
)

// Local file storage
type Local struct {
	WorkingDir string
	HistoryDir string
	db         *database.DB
}

// NewLocal return new LocalStorage
func NewLocal(db *database.DB) *Local {
	var local = Local{}
	local.db = db
	local.WorkingDir = DefaultWorkingDir
	local.HistoryDir = DefaultHistoryDir
	if w := os.Getenv("IMAGE_WORKING_DIR"); w != "" {
		local.WorkingDir = w
	}
	if w := os.Getenv("IMAGE_HISTORY_DIR"); w != "" {
		local.HistoryDir = w
	}
	return &local
}

// NewFile from fileheader
func (lc *Local) NewFile(reader io.Reader, fileName string) (string, error) {
	dst, err := lc.correctFileName(fileName)
	if err != nil {
		return "", err
	}

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, reader)

	// Save new file to database if this file created successfully
	fileModel := database.File{Fullname: dst}
	err = lc.db.CreateFile(&fileModel)

	// If failed to save to database, delete the file
	if err != nil {
		lc.DeleteFile(dst)
		return "", err
	}
	return dst, nil
}

// RollbackNewFile will try to rollback of action newfile
func (lc *Local) RollbackNewFile(path string) (err error) {
	lc.DeleteFile(path)
	return
}

// ReplaceFile in storage
func (lc *Local) ReplaceFile(path string, file io.Reader, hasFallback bool) (string, error) {
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
	_, err = lc.NewFile(file, filepath.Base(path))

	// If failed to create file, rollback action delete file
	if err != nil {
		lc.RollbackDeleteFile(path, bkDelFile)
		return "", err
	}
	return bkDelFile, nil
}

//RenameFile in storage
func (lc *Local) RenameFile(path, newName string) (string, error) {
	// Find the file in database
	file, err := lc.db.GetFileByName(path)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(path)
	newPath := filepath.Join(dir, newName)
	err = os.Rename(path, newPath)
	if err != nil {
		return "", err
	}

	// Save rename action to database.
	// If failed to save action, rename to the origin one
	if err := lc.db.RenameFile(file, newPath); err != nil {
		os.Rename(newPath, path)
		return "", err
	}
	return newPath, nil
}

// RollbackRenameFile will try to rollback of action renamefile
func (lc *Local) RollbackRenameFile(path, newName string) (err error) {
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
func (lc *Local) DeleteFile(fileName string) (string, error) {
	// Find the file from database, if no file found return error
	file, err := lc.db.GetFileByName(fileName)
	if err != nil {
		return "", err
	}

	// Copy the file to history zone
	dst, err := copyFile(fileName, fileName, true)
	if err != nil {
		return "", err
	}

	// Delete the actual file in working zone
	err = os.Remove(fileName)
	if err != nil {
		// Delete copyfile when removing the file is getting error
		os.Remove(dst)
		return "", err
	}

	// Remove the file in database and its delete action
	err = lc.db.DeleteFile(file, dst)

	// If can not save the file, copy the from the history zone to working zone and remove the file in history zone
	if err != nil {
		if _, cfErr := copyFile(dst, fileName, false); cfErr == nil {
			os.Remove(dst)
		}
		return "", err
	}
	return dst, nil
}

// RollbackDeleteFile will try to rollback action deletefile
// Get the backup file and copy that file to working zone
func (lc *Local) RollbackDeleteFile(fileName, bkFile string) (err error) {
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
func (lc *Local) GetResized2DImage(filename string, width, height uint) (image.Image, error) {
	var path = filepath.Join(lc.WorkingDir, filename)
	//Check if file extention is valid
	var ext = filepath.Ext(path)
	if !IsValidExt(ext) {
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
