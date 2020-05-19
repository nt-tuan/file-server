package models

import "github.com/thanhtuan260593/file-server/database"

//ImageIDReq model bind id from uri
type ImageIDReq struct {
	ID uint `uri:"id" binding:"required"`
}

//ImagesReq model bind request images model
type ImagesReq struct {
	PageSize    uint     `form:"pageSize"`
	PageCurrent uint     `form:"pageCurrent"`
	OrderBy     []string `form:"orderBy" validate:"oneof=id created_at fullname"`
	OrderDir    []string `form:"orderDir" validate:"oneof=asc desc"`
	Tags        []string `json:"tags"`
}

//ImageRenameReq bind rename request model
type ImageRenameReq struct {
	Name string `json:"name" binding:"required"`
}

//ImageNewReq bind new file request model
type ImageNewReq struct {
	Name string   `form:"name" binding:"required"`
	Tags []string `form:"tags"`
}

//ImageInfoRes model
type ImageInfoRes struct {
	Fullname string   `json:"fullname"`
	Tags     []string `json:"tags"`
}

//NewImageInfoRes model
func NewImageInfoRes(img *database.File) *ImageInfoRes {
	rs := ImageInfoRes{}
	rs.Fullname = img.Fullname
	if img.Tags != nil {
		rs.Tags = make([]string, len(img.Tags))
		for i, tag := range img.Tags {
			rs.Tags[i] = tag.ID
		}
	}
	return &rs
}
