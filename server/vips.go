package server

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
	"github.com/ptcoffee/image-server/server/models"
)

func (s *Server) bindResizeRequest(c *gin.Context) (*bimg.Image, *models.ResizeImageReq, error) {
	var model models.ResizeImageReq
	if err := c.BindUri(&model); err != nil {
		return nil, nil, err
	}
	s.config.CorrectImageModel(&model)
	buffer, err := s.storage.GetImageBuffer(model.FileName)
	if err != nil {
		return nil, nil, err
	}
	img := bimg.NewImage(buffer)
	return img, &model, nil
}

func responseImage(c *gin.Context, data []byte) {
	contentType := "image/" + bimg.NewImage(data).Type()
	extraHeaders := map[string]string{
		"Content-Disposition": `inline`,
	}
	c.DataFromReader(200, int64(len(data)), contentType, bytes.NewBuffer(data), extraHeaders)
}

// HandleResize godocs
// Id HandleResize
// @Summary Get a resized image
// @Param width path uint true "Width of image. Zero if resize scaled on its height"
// @Param height path uint true "Height of image. Zero if resize scaled on its width"
// @Param /name path string true "Image local path"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /images/size/{width}/{height}/{/name} [get]
func (s *Server) HandleResize(c *gin.Context) {
	img, model, err := s.bindResizeRequest(c)
	if err != nil {
		errorJSON(c, err)
		return
	}
	data, err := img.ResizeAndCrop(model.Width, model.Height)
	if err != nil {
		errorJSON(c, err)
	}
	responseImage(c, data)
}

// HandleGetWebpImage godocs
// Id HandleGetWebpImage
// @Summary Get a webp image
// @Param width path uint true "Width of image. Zero if resize scaled on its height"
// @Param height path uint true "Height of image. Zero if resize scaled on its width"
// @Param /name path string true "Image local path"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /images/webp/{width}/{height}/{/name} [get]
func (s *Server) HandleGetWebpImage(c *gin.Context) {
	img, model, err := s.bindResizeRequest(c)
	if err != nil {
		errorJSON(c, err)
		return
	}
	resized, err := img.ResizeAndCrop(model.Width, model.Height)
	if err != nil {
		errorJSON(c, err)
		return
	}
	data, err := bimg.NewImage(resized).Convert(bimg.WEBP)
	if err != nil {
		errorJSON(c, err)
		return
	}
	responseImage(c, data)
}
