package local

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

func tryGetNotExistFilename(path string) (rt string, err error) {
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(path, ext)
	rt = path
	for count := 1; count < MaxDuplicateFile; count++ {
		_, err = os.Stat(rt)
		if err != nil {
			if !os.IsNotExist(err) {
				break
			} else {
				return
			}
		}
		rt = fmt.Sprintf("%v_%v.%v", name, count, ext)
	}
	return
}

func (lc *Local) correctFileName(source string) (string, error) {
	//TODO: change file name to a valid one
	base := filepath.Base(source)
	dest := filepath.Clean(base)
	path := filepath.Join(lc.WorkingDir, dest)
	if !fileExists(path) {
		return path, nil
	}
	return "", ErrFileExisted
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

func IsValidExt(ext string) bool {
	for _, item := range ValidExts {
		if item == ext {
			return true
		}
	}
	return false
}

func (lc *Local) moveToTrash(filename string) (dst string, err error) {
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

	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
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
