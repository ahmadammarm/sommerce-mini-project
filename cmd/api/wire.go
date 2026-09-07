//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ahmadammarm/sommerce-mini-project/internal/alamat"
	"github.com/ahmadammarm/sommerce-mini-project/internal/category"
	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
	"github.com/ahmadammarm/sommerce-mini-project/internal/router"
	"github.com/ahmadammarm/sommerce-mini-project/internal/toko"
	"github.com/ahmadammarm/sommerce-mini-project/internal/transaction"
	"github.com/ahmadammarm/sommerce-mini-project/internal/upload"
	"github.com/ahmadammarm/sommerce-mini-project/internal/user"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/config"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/emsifa"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/uid"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

func InitializeApp() (*fiber.App, error) {
	wire.Build(
		config.NewDatabase,
		config.NewValidator,
		uid.NewIDGenerator,
		emsifa.NewEmsifaClient,
		user.UserSet,
		toko.TokoSet,
		alamat.AlamatSet,
		category.CategorySet,
		produk.ProdukSet,
		transaction.TransactionSet,
		upload.UploadSet,
		router.NewRouter,
	)
	return nil, nil
}
