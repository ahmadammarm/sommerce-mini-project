package produk

import (
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type Produk struct {
	ID             string    `gorm:"primaryKey;type:varchar(50);column:id"`
	NamaProduk     string    `gorm:"column:nama_produk;type:varchar(255)"`
	Slug           string    `gorm:"column:slug;type:varchar(255)"`
	HargaReseller  int       `gorm:"column:harga_reseller;type:int"`
	HargaKonsumen  int       `gorm:"column:harga_konsumen;type:int"`
	Stok           int       `gorm:"column:stok;type:int"`
	Deskripsi      string    `gorm:"column:deskripsi;type:text"`
	IdToko         string    `gorm:"column:id_toko;type:varchar(50);not null"`
	IdCategory     string    `gorm:"column:id_category;type:varchar(50);not null"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (p *Produk) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = cuid.New()
	return
}

type FotoProduk struct {
	ID        string    `gorm:"primaryKey;type:varchar(50);column:id"`
	IdProduk  string    `gorm:"column:id_produk;type:varchar(50);not null"`
	Url       string    `gorm:"column:url;type:varchar(255)"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (f *FotoProduk) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = cuid.New()
	return
}

type LogProduk struct {
	ID             string    `gorm:"primaryKey;type:varchar(50);column:id"`
	IdProduk       string    `gorm:"column:id_produk;type:varchar(50);not null"`
	NamaProduk     string    `gorm:"column:nama_produk;type:varchar(255)"`
	Slug           string    `gorm:"column:slug;type:varchar(255)"`
	HargaReseller  int       `gorm:"column:harga_reseller;type:int"`
	HargaKonsumen  int       `gorm:"column:harga_konsumen;type:int"`
	Deskripsi      string    `gorm:"column:deskripsi;type:text"`
	IdToko         string    `gorm:"column:id_toko;type:varchar(50);not null"`
	IdCategory     string    `gorm:"column:id_category;type:varchar(50);not null"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (l *LogProduk) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = cuid.New()
	return
}
