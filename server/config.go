package server

import (
	"os"
	"strconv"
)

//Config of server
type Config struct {
	MaxWidth  uint
	MaxHeight uint
}

//NewConfig instance
func NewConfig() *Config {
	var config = Config{MaxWidth, MaxHeight}
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
