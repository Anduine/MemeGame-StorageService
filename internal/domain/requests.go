package domain

type DeleteImagesRequest struct {
	Filenames []string `json:"filenames"`
}

type DeleteAvatarRequest struct {
	Filename string `json:"filename"`
}
