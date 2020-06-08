package server

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/thanhtuan260593/file-server/server/models"
	localstorage "github.com/thanhtuan260593/file-server/storages/local"
)

// func parseImageToReader(lc *localstorage.Storage, img image.Image, ext string) (io.Reader, int64, error) {
// 	if !lc.IsValidExt(ext) {
// 		return nil, 0, localstorage.ErrFileExtInvalid
// 	}
// 	var resizedBuffer bytes.Buffer
// 	var mybytes []byte
// 	var err error
// 	switch ext {
// 	case localstorage.PngExt:
// 		var encoder = png.Encoder{
// 			CompressionLevel: png.BestCompression,
// 		}
// 		if err := encoder.Encode(&resizedBuffer, img); err != nil {
// 			return nil, 0, err
// 		}
// 		mybytes = resizedBuffer.Bytes()
// 		mybytes, err = pngquant.CompressBytes(mybytes, "1")
// 		if err != nil {
// 			return nil, 0, err
// 		}
// 	default:
// 		return nil, 0, localstorage.ErrFileExtInvalid
// 	}

// 	reader := bytes.NewReader(mybytes)
// 	return reader, int64(resizedBuffer.Len()), nil
// }

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
