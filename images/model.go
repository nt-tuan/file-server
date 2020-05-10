package images

//ImageReq model
type ImageReq struct {
	Width    uint   `uri:"width" binding:"required"`
	Height   uint   `uri:"height" binding:"required"`
	FileName string `uri:"name" binding:"required"`
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
