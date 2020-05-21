package localstorage

import "errors"

//DefaultWorkingDir global value
var DefaultWorkingDir string = "/files/images"

//DefaultHistoryDir global value
var DefaultHistoryDir string = "/files/_history"

//ServerImageURL value
var ServerImageURL string

//Expected errors
var (
	ErrFileNotFound   = errors.New("file-not-found")
	ErrFileNotRead    = errors.New("file-not-read")
	ErrFileExtInvalid = errors.New("file-ext-invalid")
	ErrFileExisted    = errors.New("file-existed")
)

//ValidNameChars is collection of accepted characters
var ValidNameChars = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM-"

//MaxDuplicateFile value
var MaxDuplicateFile = 2020

//PngExt, SvgExt is extensions
var (
	PngExt = ".png"
	SvgExt = ".svg"
)
