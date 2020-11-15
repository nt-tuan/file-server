package server

import (
	"os"
	"strconv"

	"github.com/ptcoffee/image-server/server/models"
)

//MaxBound of image
var (
	DefaultMaxWidth  int = 4000
	DefaultMaxHeight int = 2000
)

//Config of server
type Config struct {
	MaxWidth  int
	MaxHeight int
}

//NewConfig instance
func NewConfig() *Config {
	var config = Config{DefaultMaxWidth, DefaultMaxHeight}
	maxWidth := os.Getenv("IMAGE_MAX_WIDTH")
	if w, err := strconv.ParseInt(maxWidth, 10, 32); err == nil {
		config.MaxWidth = int(w)
	}

	maxHeight := os.Getenv("IMAGE_MAX_HEIGHT")
	if h, err := strconv.ParseInt(maxHeight, 10, 32); err == nil {
		config.MaxHeight = int(h)
	}
	return &config
}

//CorrectImageModel image request parameters
func (conf *Config) CorrectImageModel(img *models.ResizeImageReq) {
	if img.Width > conf.MaxWidth {
		img.Width = conf.MaxWidth
	}
	if img.Height > conf.MaxHeight {
		img.Height = conf.MaxHeight
	}
}
