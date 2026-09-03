package category

type CategoryRequest struct {
	NamaCategory string `json:"nama_category" validate:"required"`
}

type CategoryResponse struct {
	ID           string `json:"id"`
	NamaCategory string `json:"nama_category"`
}
