package imaging

import "errors"

var tempPath = "files/_temp"
var (
	//ErrExtNotSupported error
	ErrExtNotSupported error = errors.New("extension-not-supported")
)

var supportedExts = []string{".png"}

func checkExtension(checkingExt string) (err error) {
	err = ErrExtNotSupported
	for _, ext := range supportedExts {
		if ext == checkingExt {
			err = nil
			return
		}
	}
	return
}
