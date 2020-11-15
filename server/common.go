package server

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/ptcoffee/image-server/server/models"
	localstorage "github.com/ptcoffee/image-server/storages/local"
)

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
		var model = models.ErrorRes{}
		model.Err = err.Error()
		c.AbortWithStatusJSON(400, &model)
	}
	return err
}
