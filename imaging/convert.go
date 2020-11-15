package imaging

import (
	"bytes"
	"image"
	"io"

	"github.com/nickalie/go-webpbin"
)

// ConvertToWebp will return a webp format of image
func ConvertToWebp(img image.Image) (io.Reader, int, error) {
	var buf bytes.Buffer
	if err := webpbin.Encode(&buf, img); err != nil {
		return nil, 0, err
	}

	return &buf, buf.Len(), nil
}
