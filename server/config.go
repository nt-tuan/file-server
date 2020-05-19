package server

import (
	"os"
	"strconv"

	"github.com/thanhtuan260593/file-server/server/models"
)

//MaxBound of image
var (
	DefaultMaxWidth  uint = 4000
	DefaultMaxHeight uint = 2000
)

//Config of server
type Config struct {
	MaxWidth  uint
	MaxHeight uint
}

//NewConfig instance
func NewConfig() *Config {
	var config = Config{DefaultMaxWidth, DefaultMaxHeight}
	maxWidth := os.Getenv("IMAGE_MAX_WIDTH")
	if w, err := strconv.ParseUint(maxWidth, 10, 32); err == nil {
		config.MaxWidth = uint(w)
	}

	maxHeight := os.Getenv("IMAGE_MAX_HEIGHT")
	if h, err := strconv.ParseUint(maxHeight, 10, 32); err == nil {
		config.MaxHeight = uint(h)
	}
	return &config
}

//CorrectImageModel image request parameters
func (conf *Config) CorrectImageModel(img *models.ImageFileReq) {
	if img.Width > conf.MaxWidth {
		img.Width = conf.MaxWidth
	}
	if img.Height > conf.MaxHeight {
		img.Height = conf.MaxHeight
	}
}
