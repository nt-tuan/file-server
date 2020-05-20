package server

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thanhtuan260593/file-server/server/models"
)

//HandleUploadImage response the image url if success
func (s *Server) HandleUploadImage(c *gin.Context) {
	var model models.ImageNewReq
	if err := c.Bind(&model); err != nil {
		return
	}
	// single file
	reader, err := getFileFromGinContext(c)
	if err != nil {
		errorJSON(c, err)
		return
	}
	if _, err := s.storage.AddFile(reader, model.Name); err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(200, nil)
}

//HandleResize image
func (s *Server) HandleResize(c *gin.Context) {
	var model models.ImageFileReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}
	s.config.CorrectImageModel(&model)
	img, err := s.storage.GetResized2DImage(model.FileName, model.Width, model.Height)
	if err != nil {
		errorJSON(c, err)
		return
	}
	var ext = filepath.Ext(model.FileName)
	resReader, contentLength, err := parseImageToReader(s.storage, img, ext)
	if err != nil {
		errorJSON(c, err)
		return
	}

	contentType := "image/" + strings.Trim(ext, ".")
	extraHeaders := map[string]string{
		"Content-Disposition": `inline`,
	}
	c.DataFromReader(200, contentLength, contentType, resReader, extraHeaders)
}

// HandleDeleteImage will delete image in storage
func (s *Server) HandleDeleteImage(c *gin.Context) {
	var model models.ImageIDReq
	if err := c.BindUri(&model); err != nil {
		return
	}
	file, err := s.db.GetFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
	}
	if _, err := s.storage.DeleteFile(file.Fullname); err != nil {
		errorJSON(c, err)
		return
	}
}

// HandleRenameImage will rename image
func (s *Server) HandleRenameImage(c *gin.Context) {
	var model models.ImageRenameReq
	var modelID models.ImageIDReq
	if err := errorJSON(c, c.BindJSON(&model)); err != nil {
		return
	}
	if err := errorJSON(c, c.BindUri(&modelID)); err != nil {
		return
	}
	file, err := s.db.GetFileByID(modelID.ID)
	if err != nil {
		errorJSON(c, err)
		return
	}
	if _, err := s.storage.RenameFile(file.Fullname, model.Name); err != nil {
		errorJSON(c, err)
		return
	}
}

//HandleReplaceImage will replace image
func (s *Server) HandleReplaceImage(c *gin.Context) {
	var model models.ImageIDReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}

	file, err := s.db.GetFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
		return
	}

	reader, err := getFileFromGinContext(c)
	if err != nil {
		errorJSON(c, err)
		return
	}
	s.storage.ReplaceFile(file.Fullname, reader)
	c.JSON(200, nil)
}

//HandleGetImages by tags
func (s *Server) HandleGetImages(c *gin.Context) {
	var model models.ImagesReq
	if err := errorJSON(c, c.BindQuery(&model)); err != nil {
		return
	}
	var orders []string
	if model.OrderBy != nil {
		orders = make([]string, len(model.OrderBy))
		for i, by := range model.OrderBy {
			dir := "asc"
			if model.OrderDir != nil && len(model.OrderDir) > i {
				dir = model.OrderDir[i]
			}
			orders[i] = fmt.Sprintf("%v %v", by, dir)
		}
	}

	imgs, err := s.db.GetFiles(model.Tags, model.PageCurrent, model.PageSize, orders)
	if err != nil {
		errorJSON(c, err)
		return
	}
	var rs []*models.ImageInfoRes
	rs = make([]*models.ImageInfoRes, len(imgs))
	for i, img := range imgs {
		rs[i] = models.NewImageInfoRes(&img)
	}
	c.JSON(200, rs)
}

//HandleGetImageByID will response image file if exists
func (s *Server) HandleGetImageByID(c *gin.Context) {
	var model models.ImageIDReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}
	file, err := s.db.GetFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(200, file)
}

//HandleAddImageTag will response true if success
func (s *Server) HandleAddImageTag(c *gin.Context) {
	var model models.ImageTagReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}
	file, err := s.db.GetFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
		return
	}
	if err := s.db.AddTag(file, model.Tag); err != nil {
		errorJSON(c, err)
		return
	}
}

//HandleRemoveImageTag will response false if failed
func (s *Server) HandleRemoveImageTag(c *gin.Context) {
	var model models.ImageTagReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}
	file, err := s.db.GetFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
		return
	}
	if err := s.db.RemoveTag(file, model.Tag); err != nil {
		errorJSON(c, err)
		return
	}
}
