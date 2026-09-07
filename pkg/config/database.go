package config

import (
	"os"

	"github.com/ahmadammarm/sommerce-mini-project/internal/alamat"
	"github.com/ahmadammarm/sommerce-mini-project/internal/category"
	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
	"github.com/ahmadammarm/sommerce-mini-project/internal/toko"
	"github.com/ahmadammarm/sommerce-mini-project/internal/transaction"
	"github.com/ahmadammarm/sommerce-mini-project/internal/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Konfigurasi Laragon: user 'root', password kosong
		dsn = "root:@tcp(127.0.0.1:3306)/sommerce?charset=utf8mb4&parseTime=True&loc=Local"
	}
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Otomatis membuat tabel-tabel di MySQL jika belum ada
	err = db.AutoMigrate(
		&user.User{},
		&toko.Toko{},
		&alamat.Alamat{},
		&category.Category{},
		&produk.Produk{},
		&produk.LogProduk{},
		&transaction.Trx{},
		&transaction.DetailTrx{},
	)
	
	if err != nil {
		return nil, err
	}

	return db, nil
}
