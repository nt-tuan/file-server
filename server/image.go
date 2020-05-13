package server

import (
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thanhtuan260593/file-server/storages/local"
)

//HandleUploadImage response the image url if success
func (s *Server) HandleUploadImage(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	if file == nil {
		c.AbortWithStatusJSON(400, local.ErrFileNotFound)
		return
	}

	c.JSON(200, nil)
}

//HandlerResize image
func (s *Server) HandlerResize(c *gin.Context) {
	var model ImageReq
	if err := c.BindUri(&model); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}
	model.SelfCorrect()
	img, err := s.storage.GetResizedImage(model.FileName, model.Width, model.Height)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	var ext = filepath.Ext(model.FileName)
	resReader, contentLength, err := parseImageToReader(img, ext)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	contentType := "image/" + strings.Trim(ext, ".")
	extraHeaders := map[string]string{
		"Content-Disposition": `inline`,
	}
	c.DataFromReader(200, contentLength, contentType, resReader, extraHeaders)
}

// HandlerDeleteImage will delete image in storage
func (s *Server) HandlerDeleteImage(c *gin.Context) {

}

// HandlerRenameImage will rename image
func (s *Server) HandlerRenameImage(c *gin.Context) {

}

//HandlerReplaceImage will replace image
func (s *Server) HandlerReplaceImage(c *gin.Context) {

}

//HandlerGetImages by tags
func (s *Server) HandlerGetImages(c *gin.Context) {

}

//HandlerGetImageByID will response image file if exists
func (s *Server) HandlerGetImageByID(c *gin.Context) {

}
