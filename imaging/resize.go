package imaging

import (
	"bytes"
	"image"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// Resize return a resized image
func Resize(img image.Image, width, height uint) image.Image {
	resized := imaging.Resize(img, int(width), int(height), imaging.Lanczos)
	return resized
}

// ResizeAndEncode return reader if no errors
func ResizeAndEncode(img image.Image, ext string, width, height uint) (io.Reader, int64, error) {
	resized := Resize(img, width, height)
	return EncodeImageToReader(resized, ext)
}

func getImageReader(filename string) (io.Reader, uint64, error) {
	path := getFilePath(filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}
	dataReader := bytes.NewReader(data)
	return dataReader, uint64(len(data)), err
}

func getFilePath(filename string) string {
	return filepath.Join(tempPath, filename)
}
