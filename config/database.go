package config

import (
	"fmt"
	"time"

	"github.com/ahmadammarm/sommerce-mini-project/internal/alamat"
	"github.com/ahmadammarm/sommerce-mini-project/internal/category"
	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
	"github.com/ahmadammarm/sommerce-mini-project/internal/toko"
	"github.com/ahmadammarm/sommerce-mini-project/internal/transaction"
	"github.com/ahmadammarm/sommerce-mini-project/internal/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDatabase establishes a connection to MySQL, configures the connection pool,
// and runs the GORM AutoMigrate function.
func InitDatabase(env *EnvConfig) *gorm.DB {
	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.DBUser,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	// Connect using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// ---------------------------------------------------------
	// CONNECTION POOLING CONFIGURATION
	// ---------------------------------------------------------
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to retrieve underlying SQL DB: " + err.Error())
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	// ---------------------------------------------------------
	// AUTO MIGRATE
	// ---------------------------------------------------------
	err = db.AutoMigrate(
		&user.User{},
		&alamat.Alamat{},
		&toko.Toko{},
		&category.Category{},
		&produk.Produk{},
		&produk.FotoProduk{},
		&produk.LogProduk{},
		&transaction.Trx{},
		&transaction.DetailTrx{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	return db
}
