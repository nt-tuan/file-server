package images

//ImageReq model
type ImageReq struct {
	Width    uint   `uri:"width"`
	Height   uint   `uri:"height"`
	FileName string `uri:"name" binding:"required"`
}

//UploadReq model
type UploadReq struct {
	Dir string `uri:"dir"`
}

//SelfCorrect image request parameters
func (img *ImageReq) SelfCorrect() {
	if img.Width > MaxWidth {
		img.Width = MaxWidth
	}
	if img.Height > MaxHeight {
		img.Height = MaxHeight
	}
}

//ImageRes is model returned to client
type ImageRes struct {
	URL uint `form:"url"`
}
