package images

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func changeFileName(source string, subdir string) (string, error) {
	//TODO: change file name to a valid one
	dest := filepath.Clean(source)
	ext := filepath.Ext(dest)
	filename := strings.TrimSuffix(dest, ext)
	dirPath := filepath.Join(LocalImagePath, subdir)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", err
	}
	path := filepath.Join(LocalImagePath, subdir, filename+ext)
	if !fileExists(path) {
		return path, nil
	}
	return "", ErrFileExists
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

func isValidExt(ext string) bool {
	for _, item := range ValidExts {
		if item == ext {
			return true
		}
	}
	return false
}

func parseImageToReader(img image.Image, ext string) (io.Reader, int64, error) {
	if !isValidExt(ext) {
		return nil, 0, ErrFileExtInvalid
	}
	var resizedBuffer bytes.Buffer
	switch ext {
	case PngExt:
		if err := png.Encode(&resizedBuffer, img); err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, ErrFileExtInvalid
	}

	reader := bytes.NewReader(resizedBuffer.Bytes())
	return reader, int64(resizedBuffer.Len()), nil
}
