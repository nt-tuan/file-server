package main

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thanhtuan260593/file-server/images"
)

func init() {
	images.LocalImagePath = "./files/images/" //os.Getenv("LOCAL_IMAGE_PATH")
	images.ServerImageURL = "/images"
	if v, err := strconv.ParseUint(os.Getenv("IMAGE_MAX_WIDTH"), 10, 32); err == nil {
		images.MaxWidth = uint(v)
	}
	if v, err := strconv.ParseUint(os.Getenv("IMAGE_MAX_HEIGHT"), 10, 32); err == nil {
		images.MaxHeight = uint(v)
	}
}

//Main func
func main() {
	port := os.Getenv("PORT")
	router := gin.Default()
	imageGroup := router.Group(images.ServerImageURL)
	imageGroup.Static("/static", images.LocalImagePath)
	imageGroup.GET("/size/:width/:height/:name", images.HandleResizeImage)
	imageGroup.POST("/upload", images.HandleUploadImage)

	// Listen and serve on port
	router.Run(port)
}
