package imaging

import (
	"bytes"
	"image"
	"image/png"
	"io"

	pngquant "github.com/yusukebe/go-pngquant"
)

// CompressImage return reader if no error
func CompressImage(img image.Image) (io.Reader, int64, error) {
	var resizedBuffer bytes.Buffer
	var mybytes []byte
	var err error
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
	reader := bytes.NewReader(mybytes)
	return reader, int64(len(mybytes)), nil
}
