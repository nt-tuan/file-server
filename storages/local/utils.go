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
)

//GetPhysicalWorkingPath from client path
func (lc *Storage) GetPhysicalWorkingPath(clientPath string) (rs string) {
	rs = filepath.Join(lc.WorkingDir, clientPath)
	return
}

//GetPhysicalHistoricalPath from client path
func (lc *Storage) GetPhysicalHistoricalPath(clientPath string) (rs string) {
	rs = filepath.Join(lc.HistoryDir, clientPath)
	return
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getImageFromPath(filepath string) (image.Image, error) {
	if !fileExists(filepath) {
		return nil, ErrFileNotFound
	}
	data, err := ioutil.ReadFile(filepath)
	dataReader := bytes.NewReader(data)
	imageData, _, err := image.Decode(dataReader)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

//IsValidExt return true if file extension is a valid extension
func (lc *Storage) IsValidExt(ext string) bool {
	for _, item := range lc.ValidExts {
		if item == ext {
			return true
		}
	}
	return false
}

func (lc *Storage) moveToTrash(filename string) (dst string, err error) {
	base := filepath.Base(filename)
	dst, err = copyFile(filename, filepath.Join(lc.HistoryDir, base), true)
	return
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func copyFile(src, dst string, forceIfDestExisted bool) (dest string, err error) {
	sfi, err := os.Stat(src)
	dest = dst
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		err = fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
		return
	}

	//
	_, err = os.Stat(dst)
	if err == nil {
		if forceIfDestExisted {
			var destErr error
			dest, destErr = tryGetNotExistFilename(dst)
			if destErr != nil {
				return "", destErr
			}
		} else {
			return "", ErrFileExisted
		}
	}

	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
	}

	if err = os.Link(src, dest); err == nil {
		return
	}
	err = copyFileContents(src, dest)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

//RemoveContents in dir
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
