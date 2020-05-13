package images

import (
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

//HandleUploadImage response the image url if success
func HandleUploadImage(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	var model UploadReq
	if err := c.BindUri(&model); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}
	if file == nil {
		c.AbortWithStatusJSON(400, ErrFileNotFound.Error())
		return
	}
	path, err := changeFileName(file.Filename, model.Dir)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}
	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	rsfile := filepath.Base(path)
	c.JSON(200, rsfile)
}

//HandleResizeImage response resized image to client
func HandleResizeImage(c *gin.Context) {
	var model ImageReq
	if err := c.BindUri(&model); err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}
	model.SelfCorrect()

	var path = filepath.Join(LocalImagePath, model.FileName)
	//Check if file extention is valid
	var ext = filepath.Ext(path)
	if !isValidExt(ext) {
		c.AbortWithStatusJSON(400, ErrFileExtInvalid)
		return
	}
	imageData, err := getImageFromPath(path)
	if err != nil {
		c.AbortWithStatusJSON(400, err.Error())
		return
	}

	//Resize image
	resized := imaging.Resize(imageData, int(model.Width), int(model.Height), imaging.Lanczos)
	resReader, contentLength, err := parseImageToReader(resized, ext)
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
