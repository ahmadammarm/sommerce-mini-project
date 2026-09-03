package produk

type CreateProdukRequest struct {
	NamaProduk    string   `json:"nama_produk" validate:"required"`
	HargaReseller int      `json:"harga_reseller" validate:"required,min=0"`
	HargaKonsumen int      `json:"harga_konsumen" validate:"required,min=0"`
	Stok          int      `json:"stok" validate:"required,min=0"`
	Deskripsi     string   `json:"deskripsi" validate:"required"`
	IdCategory    string   `json:"id_category" validate:"required"`
	Photos        []string `json:"photos"` // Array of URLs from upload API
}

type ProductFilterDTO struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Search     string `query:"search"`
	CategoryId string `query:"category_id"`
}

type ProdukResponse struct {
	ID            string   `json:"id"`
	NamaProduk    string   `json:"nama_produk"`
	Slug          string   `json:"slug"`
	HargaReseller int      `json:"harga_reseller"`
	HargaKonsumen int      `json:"harga_konsumen"`
	Stok          int      `json:"stok"`
	Deskripsi     string   `json:"deskripsi"`
	IdCategory    string   `json:"id_category"`
	IdToko        string   `json:"id_toko"`
	Photos        []string `json:"photos"`
}
