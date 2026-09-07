package transaction

import (
	"time"

	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
)

type Trx struct {
	ID               string    `gorm:"primaryKey;type:varchar(50);column:id"`
	IdUser           string    `gorm:"column:id_user;type:varchar(50);not null;index"`
	AlamatPengiriman string    `gorm:"column:alamat_pengiriman;type:varchar(50);not null"`
	HargaTotal       int       `gorm:"column:harga_total;type:int"`
	IdempotencyKey   string    `gorm:"column:idempotency_key;type:varchar(100);unique;not null"`
	KodeInvoice      string    `gorm:"column:kode_invoice;type:varchar(255);unique"`
	MethodBayar      string    `gorm:"column:method_bayar;type:varchar(255)"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime"`
}

type DetailTrx struct {
	ID          string           `gorm:"primaryKey;type:varchar(50);column:id"`
	IdTrx       string           `gorm:"column:id_trx;type:varchar(50);not null;index"`
	IdLogProduk string           `gorm:"column:id_log_produk;type:varchar(50);not null;index"`
	LogProduk   produk.LogProduk `gorm:"foreignKey:IdLogProduk"`
	IdToko      string           `gorm:"column:id_toko;type:varchar(50);not null;index"`
	Kuantitas   int              `gorm:"column:kuantitas;type:int"`
	HargaTotal  int              `gorm:"column:harga_total;type:int"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt   time.Time        `gorm:"column:created_at;autoCreateTime"`
}
