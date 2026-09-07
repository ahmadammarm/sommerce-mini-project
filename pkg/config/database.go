package config

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Standard MySQL DSN
		dsn = "root:root@tcp(127.0.0.1:3306)/sommerce?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
