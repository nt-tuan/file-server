package localstorage

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"
	"github.com/twinj/uuid"
)

//getPath from client path
func (lc *Storage) getPath(fullname string) (rs string) {
	rs = filepath.Join(lc.WorkingDir, fullname)
	return
}

//getBackupPath from client path
func (lc *Storage) getBackupPath(clientPath string) (rs string) {
	rs = filepath.Join(lc.HistoryDir, clientPath)
	return
}

func (lc *Storage) newBackupFullname(fullname string) (rs string) {
	ext := filepath.Ext(fullname)
	name := uuid.NewV4().String()
	return name + ext
}

func tryGetNotExistFilename(path string) (string, error) {
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(path, ext)
	rt := path
	for count := 1; count < MaxDuplicateFile; count++ {
		if !fileExists(rt) {
			return rt, nil
		}
		rt = fmt.Sprintf("%v_%v%v", name, count, ext)
	}
	return "", ErrFileExisted
}

func (lc *Storage) correctFileName(source string) (string, string, error) {
	//TODO: change file name to a valid one
	clientPath := filepath.Clean(source)
	serverPath := filepath.Join(lc.WorkingDir, clientPath)
	if !fileExists(serverPath) {
		return serverPath, clientPath, nil
	}
	return "", "", ErrFileExisted
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getImageBuffer(filepath string) ([]byte, error) {
	if !fileExists(filepath) {
		return nil, ErrFileNotFound
	}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getImageFromPath(filepath string) (image.Image, error) {
	data, err := getImageBuffer(filepath)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	imageData, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

//IsValidExt return true if file extension is a valid extension
func (lc *Storage) IsValidExt(ext string) bool {
	for _, item := range lc.ValidExts {
		if strings.ToLower(item) == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

// GetFileSize return file storage size
func (lc *Storage) GetFileSize(fullname string) (int64, error) {
	path := lc.getPath(fullname)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// GetImageSize return width and height of image
func (lc *Storage) GetImageSize(fullname string) (bimg.ImageSize, error) {
	path := lc.getPath(fullname)
	data, err := bimg.Read(path)
	if err != nil {
		return bimg.ImageSize{}, err
	}
	return bimg.Size(data)
}

// RemoveContents will clear directory
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

/*
   GoLang: os.Rename() give error "invalid cross-device link" for Docker container with Volumes.
   MoveFile(source, destination) will work moving file between folders
*/

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
