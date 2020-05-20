package models

//ImageFileReq model
type ImageFileReq struct {
	Width    uint   `uri:"width"`
	Height   uint   `uri:"height"`
	FileName string `uri:"name" binding:"required"`
}
