package server

import (
	"bytes"
	"image"
	"image/png"
	"io"

	"github.com/gin-gonic/gin"
	localstorage "github.com/thanhtuan260593/file-server/storages/local"
)

func parseImageToReader(lc *localstorage.Storage, img image.Image, ext string) (io.Reader, int64, error) {
	if !lc.IsValidExt(ext) {
		return nil, 0, localstorage.ErrFileExtInvalid
	}
	var resizedBuffer bytes.Buffer
	switch ext {
	case localstorage.PngExt:
		if err := png.Encode(&resizedBuffer, img); err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, localstorage.ErrFileExtInvalid
	}

	reader := bytes.NewReader(resizedBuffer.Bytes())
	return reader, int64(resizedBuffer.Len()), nil
}

func getFileFromGinContext(c *gin.Context) (io.Reader, error) {
	fileHeader, _ := c.FormFile("file")
	if fileHeader == nil {
		return nil, localstorage.ErrFileNotFound
	}
	reader, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func errorJSON(c *gin.Context, err error) error {
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
	}
	return err
}
