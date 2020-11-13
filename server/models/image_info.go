package models

import "github.com/ptcoffee/image-server/database"

//IDReq model bind id from uri
type IDReq struct {
	ID uint `uri:"id" binding:"required"`
}

//ImagesReq model bind request images model
type ImagesReq struct {
	PageSize    uint     `form:"pageSize"`
	PageCurrent uint     `form:"pageCurrent"`
	OrderBy     []string `form:"orderBy"`
	OrderDir    []string `form:"orderDir"`
	Tags        []string `form:"tags"`
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
	ID       uint     `json:"id"`
	Fullname string   `json:"fullname"`
	Tags     []string `json:"tags"`
}

//NewImageInfoRes model
func NewImageInfoRes(img *database.File) *ImageInfoRes {
	rs := ImageInfoRes{}
	rs.Fullname = img.Fullname
	rs.ID = img.ID
	if img.Tags != nil {
		rs.Tags = make([]string, len(img.Tags))
		for i, tag := range img.Tags {
			rs.Tags[i] = tag.ID
		}
	}
	return &rs
}
