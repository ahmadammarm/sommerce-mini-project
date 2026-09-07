package produk

import (
	"time"
)

type Produk struct {
	ID            string       `gorm:"primaryKey;type:varchar(50);column:id"`
	NamaProduk    string       `gorm:"column:nama_produk;type:varchar(255);index"`
	Slug          string       `gorm:"column:slug;type:varchar(255)"`
	HargaReseller int          `gorm:"column:harga_reseller;type:int"`
	HargaKonsumen int          `gorm:"column:harga_konsumen;type:int"`
	Stok          int          `gorm:"column:stok;type:int"`
	Deskripsi     string       `gorm:"column:deskripsi;type:text"`
	IdToko        string       `gorm:"column:id_toko;type:varchar(50);not null;index"`
	IdCategory    string       `gorm:"column:id_category;type:varchar(50);not null;index"`
	Fotos         []FotoProduk `gorm:"foreignKey:IdProduk"`
	UpdatedAt     time.Time    `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt     time.Time    `gorm:"column:created_at;autoCreateTime"`
}

type FotoProduk struct {
	ID        string    `gorm:"primaryKey;type:varchar(50);column:id"`
	IdProduk  string    `gorm:"column:id_produk;type:varchar(50);not null;index"`
	Url       string    `gorm:"column:url;type:varchar(255)"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type LogProduk struct {
	ID            string    `gorm:"primaryKey;type:varchar(50);column:id"`
	IdProduk      string    `gorm:"column:id_produk;type:varchar(50);not null;index"`
	NamaProduk    string    `gorm:"column:nama_produk;type:varchar(255)"`
	Slug          string    `gorm:"column:slug;type:varchar(255)"`
	HargaReseller int       `gorm:"column:harga_reseller;type:int"`
	HargaKonsumen int       `gorm:"column:harga_konsumen;type:int"`
	Deskripsi     string    `gorm:"column:deskripsi;type:text"`
	IdToko        string    `gorm:"column:id_toko;type:varchar(50);not null;index"`
	IdCategory    string    `gorm:"column:id_category;type:varchar(50);not null;index"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}
