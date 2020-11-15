package models

//ResizeImageReq is request model of resizing image api
type ResizeImageReq struct {
	Width    int    `uri:"width"`
	Height   int    `uri:"height"`
	FileName string `uri:"name" binding:"required"`
}

// ConvertImageReq is request model of converting image api
type ConvertImageReq struct {
	FileName string `uri:"name" binding:"required"`
}
