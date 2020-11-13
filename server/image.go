package server

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ptcoffee/image-server/imaging"
	"github.com/ptcoffee/image-server/server/models"
)

func getImageURL(fullname string) string {
	return "https://" + os.Getenv("BASE_PATH") + "/images/static/" + fullname
}

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
	var model models.IDReq
	if err := c.BindUri(&model); err != nil {
		return
	}
	file, err := s.db.GetFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
	}
	if err := s.storage.DeleteFile(file); err != nil {
		errorJSON(c, err)
		return
	}
	if err := s.cloudflareAPI.PurgeCache(getImageURL(file.Fullname)); err != nil {
		log.Println(err)
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
	var modelID models.IDReq
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
	if err := s.storage.RenameFile(file, model.Name); err != nil {
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
	var model models.IDReq
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
	if _, err := s.storage.ReplaceFile(file, reader); err != nil {
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
	var model models.IDReq
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

// HandlePurgeCDNCache godocs
// @Id HandlePurgeCDNCache
// @Summary Clear cache of image
// @Param id path uint true "ID of image"
// @Success 200
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image/{id}/purgeCache [post]
func (s *Server) HandlePurgeCDNCache(c *gin.Context) {
	var model models.IDReq
	if err := c.BindUri(&model); err != nil {
		return
	}
	file, err := s.db.GetFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
	}
	if err := s.cloudflareAPI.PurgeCache(getImageURL(file.Fullname)); err != nil {
		errorJSON(c, err)
	}
	c.Status(200)
}

// HandleGetImageHistory godocs
// @Id HandleGetImageHistory
// @Summary Get list of history changes of an image
// @Description Get list of images information
// @Produce  json
// @Param id path uint true "ID of image"
// @Success 200 {object} []database.FileHistory
// @Failure 400 {object} models.ErrorRes
// @Router /admin/image/{id}/history [get]
func (s *Server) HandleGetImageHistory(c *gin.Context) {
	var model models.IDReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}
	records, err := s.db.GetFileHistoryRecords(model.ID)
	if err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(200, models.NewHistoryInfos(records))
}

// HandleGetDeletedFiles godocs
// @Id HandleGetDeletedFiles
// @Summary Get list of deleted files
// @Success 200 {object} []database.FileHistory
// @Failure 400 {object} models.ErrorRes
// @Router /admin/deletedImages
func (s *Server) HandleGetDeletedFiles(c *gin.Context) {
	deletedFiles, err := s.db.GetDeletedFiles()
	if err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(200, models.NewHistoryInfos(deletedFiles))
}

// HandleRecoverDeletedFile godocs
// @Id HandleRecoverDeletedFile
// @Summary Recover a deleted file
// @Success 200 {object} models.ImageInfoRes
// @Failure 400 {object} models.ErrorRes
// @Router /admin/deletedImage/{id}/restore
func (s *Server) HandleRecoverDeletedFile(c *gin.Context) {
	var model models.IDReq
	if err := errorJSON(c, c.BindUri(&model)); err != nil {
		return
	}
	deletedFile, err := s.db.GetDeletedFileByID(model.ID)
	if err != nil {
		errorJSON(c, err)
		return
	}
	if deletedFile.BackupFullname == nil {
		errorJSON(c, errors.New("file can not restored"))
	}
	restoredFile, err := s.storage.RestoreDeletedFile(*deletedFile)
	if err != nil {
		errorJSON(c, err)
		return
	}
	c.JSON(200, models.NewImageInfoRes(restoredFile))
}
