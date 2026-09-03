package toko

type UpdateTokoRequest struct {
	NamaToko string `json:"nama_toko" validate:"required,min=3"`
	UrlFoto  string `json:"url_foto"` // Will be populated by the file upload API
}

type TokoResponse struct {
	ID       string `json:"id"`
	NamaToko string `json:"nama_toko"`
	UrlFoto  string `json:"url_foto"`
}
