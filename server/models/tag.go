package models

//ImageTagReq model for add/remove tag from image
type ImageTagReq struct {
	ID  uint   `uri:"id" binding:"required"`
	Tag string `uri:"tag" binding:"required"`
}
