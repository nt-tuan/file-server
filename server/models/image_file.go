package models

//ImageFileReq model
type ImageFileReq struct {
	Width    uint   `uri:"width" binding:"required"`
	Height   uint   `uri:"height" binding:"required"`
	FileName string `uri:"name" binding:"required"`
}
