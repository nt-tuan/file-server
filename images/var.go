package images

import "errors"

//LocalImagePath global value
var LocalImagePath string

//ServerImageURL value
var ServerImageURL string

//ErrFileNotFound is file not exists in file server
var ErrFileNotFound = errors.New("file-not-found")

//ErrFileNotRead is error while reading the file
var ErrFileNotRead = errors.New("file-not-read")

//ErrFileExtInvalid is error when file ext is invalid
var ErrFileExtInvalid = errors.New("file-ext-invalid")

//ValidExts is collection of accepted image file extensions
var ValidExts = []string{PngExt, SvgExt}

//ValidNameChars is collection of accepted characters
var ValidNameChars = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM-"

//MaxDuplicateFile value
var MaxDuplicateFile = 2020

//MaxWidth of image
var MaxWidth uint = 4000

//MaxHeight of image
var MaxHeight uint = 2000

//PngExt, SvgExt is extensions
var (
	PngExt = ".png"
	SvgExt = ".svg"
)
