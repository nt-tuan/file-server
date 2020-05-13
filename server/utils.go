package server

import (
	"bytes"
	"image"
	"image/png"
	"io"

	"github.com/thanhtuan260593/file-server/storages/local"
)

func parseImageToReader(img image.Image, ext string) (io.Reader, int64, error) {
	if !local.IsValidExt(ext) {
		return nil, 0, local.ErrFileExtInvalid
	}
	var resizedBuffer bytes.Buffer
	switch ext {
	case local.PngExt:
		if err := png.Encode(&resizedBuffer, img); err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, local.ErrFileExtInvalid
	}

	reader := bytes.NewReader(resizedBuffer.Bytes())
	return reader, int64(resizedBuffer.Len()), nil
}
