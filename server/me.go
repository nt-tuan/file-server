package server

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ptcoffee/image-server/server/models"
)

// HandleUploadUserImage godocs
// @Id HandleUploadUserImage
// @Summary Upload an user image
// @Accept multipart/form-data
// @Param file formData file true "Upload file"
// @Param name formData string true "File name"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router / [put]
func (s *Server) HandleUploadUserImage(c *gin.Context) {
	var model models.ImageNewReq
	if err := c.Bind(&model); err != nil {
		return
	}
	unix := time.Now().Unix()
	name := model.Name + fmt.Sprintf("%d", unix)
	userName := c.GetString("User")
	fileName := filepath.Join("user", userName, name)
	// single file
	reader, err := getFileFromGinContext(c)
	if err != nil {
		errorJSON(c, err)
		return
	}
	file, err := s.storage.AddFile(reader, fileName, userName)
	if err != nil {
		errorJSON(c, err)
		return
	}
	s.db.AddTag(file, userName)
	c.JSON(200, models.NewImageInfoRes(file))
}

// HandleGetUserImages godocs
// @Id HandleGetUserImages
// @Summary Get list of images information of current user
// @Description Get list of images information of current user
// @Produce  json
// @Param model query models.ImagesReq false "query model"
// @Success 200 {array} models.ImageInfoRes
// @Failure 400 {object} models.ErrorRes
// @Router /admin/images [get]
func (s *Server) HandleGetUserImages(c *gin.Context) {
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
	userName := c.GetString("User")
	images, err := s.db.GetFilesByUser(model.Tags, model.Offset, model.Limit, orders, userName)
	if err != nil {
		errorJSON(c, err)
		return
	}
	var rs []*models.ImageInfoRes
	rs = make([]*models.ImageInfoRes, len(images))
	for i, img := range images {
		rs[i] = models.NewImageInfoRes(&img)
	}
	c.JSON(200, rs)
}
