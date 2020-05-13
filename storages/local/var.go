package local

import "errors"

//DefaultWorkingDir global value
var DefaultWorkingDir string

//DefaultHistoryDir global value
var DefaultHistoryDir string

//ServerImageURL value
var ServerImageURL string

//Expected errors
var (
	ErrFileNotFound   = errors.New("file-not-found")
	ErrFileNotRead    = errors.New("file-not-read")
	ErrFileExtInvalid = errors.New("file-ext-invalid")
	ErrFileExisted    = errors.New("file-existed")
)

//ValidExts is collection of accepted image file extensions
var ValidExts = []string{PngExt, SvgExt}

//ValidNameChars is collection of accepted characters
var ValidNameChars = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM-"

//MaxDuplicateFile value
var MaxDuplicateFile = 2020

//PngExt, SvgExt is extensions
var (
	PngExt = ".png"
	SvgExt = ".svg"
)
