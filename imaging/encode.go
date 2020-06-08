package imaging

import (
	"bytes"
	"image"
	"image/png"
	"io"

	localstorage "github.com/thanhtuan260593/file-server/storages/local"
	pngquant "github.com/yusukebe/go-pngquant"
)

// EncodeImageToReader return reader if no error
func EncodeImageToReader(img image.Image, ext string) (io.Reader, int64, error) {
	var resizedBuffer bytes.Buffer
	var mybytes []byte
	var err error
	if err = checkExtension(ext); err != nil {
		return nil, 0, err
	}
	switch ext {
	case localstorage.PngExt:
		var encoder = png.Encoder{
			CompressionLevel: png.BestSpeed,
		}
		if err := encoder.Encode(&resizedBuffer, img); err != nil {
			return nil, 0, err
		}
		mybytes = resizedBuffer.Bytes()
		mybytes, err = pngquant.CompressBytes(mybytes, "1")
		if err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, localstorage.ErrFileExtInvalid
	}

	reader := bytes.NewReader(mybytes)
	return reader, int64(len(mybytes)), nil
}
