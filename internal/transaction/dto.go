package transaction

import "time"

type CheckoutItemRequest struct {
	IdProduk  string `json:"id_produk" validate:"required"`
	Kuantitas int    `json:"kuantitas" validate:"required,min=1"`
}

type CheckoutRequest struct {
	AlamatPengiriman string                `json:"alamat_pengiriman" validate:"required"`
	MethodBayar      string                `json:"method_bayar" validate:"required"`
	Items            []CheckoutItemRequest `json:"items" validate:"required,min=1,dive"`
}

type LogProdukResponse struct {
	ID            string `json:"id"`
	IdProduk      string `json:"id_produk"`
	NamaProduk    string `json:"nama_produk"`
	HargaReseller int    `json:"harga_reseller"`
	HargaKonsumen int    `json:"harga_konsumen"`
	Deskripsi     string `json:"deskripsi"`
}

type DetailTransactionResponse struct {
	ID         string            `json:"id"`
	Kuantitas  int               `json:"kuantitas"`
	HargaTotal int               `json:"harga_total"`
	Produk     LogProdukResponse `json:"produk"`
}

type TransactionResponse struct {
	ID               string    `json:"id"`
	KodeInvoice      string    `json:"kode_invoice"`
	HargaTotal       int       `json:"harga_total"`
	MethodBayar      string    `json:"method_bayar"`
	AlamatPengiriman string    `json:"alamat_pengiriman"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type TransactionDetailResponse struct {
	TransactionResponse
	Details []DetailTransactionResponse `json:"details"`
}
