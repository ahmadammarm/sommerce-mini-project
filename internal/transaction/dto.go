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

type TransactionResponse struct {
	ID               string    `json:"id"`
	KodeInvoice      string    `json:"kode_invoice"`
	HargaTotal       int       `json:"harga_total"`
	MethodBayar      string    `json:"method_bayar"`
	AlamatPengiriman string    `json:"alamat_pengiriman"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}
