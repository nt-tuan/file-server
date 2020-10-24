package server

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thanhtuan260593/file-server/imaging"
	"github.com/thanhtuan260593/file-server/server/models"
)

// HandleUploadImage godocs
// @Id UploadImage
// @Summary Upload an image
// @Accept multipart/form-data
// @Param file formData file true "Upload file"
// @Param name formData string true "File name"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image [put]
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
	file, err := s.storage.AddFile(reader, model.Name)
	if err != nil {
		errorJSON(c, err)
		return
	}

	c.JSON(200, models.NewImageInfoRes(file))
}

// HandleResize godocs
// Id GetResizedImage
// @Summary Get a resized image
// @Param width path uint true "Width of image. Zero if resize scaled on its height"
// @Param height path uint true "Height of image. Zero if resize scaled on its width"
// @Param /name path string true "Image local path"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /images/size/{width}/{height}/{/name} [get]
func (s *Server) HandleResize(c *gin.Context) {
	var model models.ImageFileReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}
	s.config.CorrectImageModel(&model)
	img, err := s.storage.GetImage(model.FileName)
	if err != nil {
		errorJSON(c, err)
		return
	}
	var ext = filepath.Ext(model.FileName)
	//resReader, contentLength, err := imaging.ResizeAndEncode(img, ext, model.Width, model.Height)
	resReader, contentLength, err := imaging.ResizeAndEncode(img, ext, model.Width, model.Height)
	if err != nil {
		errorJSON(c, err)
		return
	}

	contentType := "image/" + strings.Trim(ext, ".")
	extraHeaders := map[string]string{
		"Content-Disposition": `inline`,
		"Cache-Control":       "public",
		"max-age":             "108000",
	}
	c.DataFromReader(200, int64(contentLength), contentType, resReader, extraHeaders)
}

// HandleDeleteImage godocs
// @Id DeleteImage
// @Summary Delete an image
// @Param id path uint true "ID of image"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image/{id} [delete]
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
	c.Status(200)
}

// HandleRenameImage godocs
// @Id RenameImage
// @Summary Rename an image
// @Accept application/json
// @Param model query models.ImageRenameReq true "query model"
// @Param id path uint true "ID of image"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image/{id}/rename [post]
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
	c.Status(200)
}

// HandleReplaceImage godoc
// @Id ReplaceImage
// @Summary Replace an image
// @Description replace and image
// @Accept multipart/form-data
// @Produce  json
// @Param id path uint true "ID of image"
// @Param file formData file true "Replaced file"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image/{id}/replace [post]
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
	if _, err := s.storage.ReplaceFile(file.Fullname, reader); err != nil {
		errorJSON(c, err)
		return
	}

	c.Status(200)
}

// HandleGetImages godocs
// @Id GetImages
// @Summary Get list of images information
// @Description Get list of images information
// @Produce  json
// @Param model query models.ImagesReq false "query model"
// @Success 200 {array} models.ImageInfoRes
// @Failure 400 {object} models.ErrorRes
// @Router /admin/images [get]
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

// HandleGetImageByID docs
// @Id GetImageByID
// @Summary Get an image information
// @Param id path uint true "ID of image"
// @Success 200 {object} models.ImageInfoRes
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image/{id} [get]
func (s *Server) HandleGetImageByID(c *gin.Context) {
	var model models.ImageIDReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}
	file, err := s.db.GetFullFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(200, models.NewImageInfoRes(file))
}

// HandleAddImageTag godocs
// @Id AddImageTag
// @Summary Add a tag to an image
// @Param id path uint true "ID of image"
// @Param tag path string true "Added tag"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image/{id}/tag/{tag} [put]
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
	c.Status(200)
}

// HandleRemoveImageTag godocs
// @Id RemoveImageTag
// @Summary Remove a tag from an image
// @Param id path uint true "ID of image"
// @Param tag path string true "Added tag"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image/{id}/tag/{tag} [delete]
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
