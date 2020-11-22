package models

import (
	"time"

	"github.com/ptcoffee/image-server/database"
)

//IDReq model bind id from uri
type IDReq struct {
	ID uint `uri:"id" binding:"required"`
}

// ImagesReq model bind request images model
type ImagesReq struct {
	Limit    *uint    `form:"limit"`
	Offset   *uint    `form:"offset"`
	OrderBy  []string `form:"orderBy"`
	OrderDir []string `form:"orderDir"`
	Tags     []string `form:"tags"`
}

// ImagesCountReq struct
type ImagesCountReq struct {
	Tags []string `form:"tags"`
}

// ImageCountRes struct
type ImageCountRes struct {
	Count int `json:"count"`
}

// ImageRenameReq bind rename request model
type ImageRenameReq struct {
	Name string `json:"name" binding:"required"`
}

//ImageNewReq bind new file request model
type ImageNewReq struct {
	Name string `form:"name" binding:"required"`
}

//ImageInfoRes model
type ImageInfoRes struct {
	ID       uint      `json:"id"`
	Fullname string    `json:"fullname"`
	Tags     []string  `json:"tags"`
	By       string    `json:"by"`
	At       time.Time `json:"at"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	DiskSize int64     `json:"diskSize"`
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
	rs.By = img.CreatedBy
	rs.At = img.CreatedAt
	rs.Width = img.Width
	rs.Height = img.Height
	rs.DiskSize = img.DiskSize
	return &rs
}
