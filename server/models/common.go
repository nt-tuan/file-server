package models

// Pagination model
type Pagination struct {
	Current uint `json:"current"`
	Size    uint `json:"size"`
}

// ImageOrder model
type ImageOrder struct {
	By        uint `json:"by" validate:"oneof=id created_at fullname"`
	Direction uint `json:"direction" validate:"oneof=asc desc"`
}

// ErrorRes model
type ErrorRes struct {
	Err string
}
