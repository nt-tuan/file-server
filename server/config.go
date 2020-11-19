package server

import (
	"errors"
	"os"
	"strconv"
)

//MaxBound of image
var (
	DefaultMaxWidth  int = 4000
	DefaultMaxHeight int = 2000
)

// ErrInvalidImageSize error
var ErrInvalidImageSize = errors.New("invalid-image-size")

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

// CheckImageSize return error if the image is invalid
func (config *Config) CheckImageSize(width int, height int) error {
	if width <= config.MaxWidth && height <= config.MaxHeight {
		return nil
	}
	return ErrInvalidImageSize
}
