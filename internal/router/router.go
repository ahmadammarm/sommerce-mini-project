package router

import (
	"github.com/ahmadammarm/sommerce-mini-project/internal/alamat"
	"github.com/ahmadammarm/sommerce-mini-project/internal/category"
	"github.com/ahmadammarm/sommerce-mini-project/internal/middleware"
	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
	"github.com/ahmadammarm/sommerce-mini-project/internal/toko"
	"github.com/ahmadammarm/sommerce-mini-project/internal/transaction"
	"github.com/ahmadammarm/sommerce-mini-project/internal/upload"
	"github.com/ahmadammarm/sommerce-mini-project/internal/user"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(
	userHandler *user.UserHandler,
	tokoHandler *toko.TokoHandler,
	alamatHandler *alamat.AlamatHandler,
	categoryHandler *category.CategoryHandler,
	produkHandler *produk.ProdukHandler,
	transactionHandler *transaction.TransactionHandler,
	uploadHandler *upload.UploadHandler,
) *fiber.App {
	app := fiber.New()

	// 1. Security Headers (Helmet)
	app.Use(middleware.SecurityHeaders())

	// 2. CORS
	app.Use(middleware.Cors())

	// 3. Global Rate Limiter (Anti DDoS)
	app.Use(middleware.GlobalRateLimiter())

	// Serve Static Files for Uploads
	app.Static("/uploads", "./public/uploads")

	api := app.Group("/api")

	// --- Upload ---
	api.Post("/upload", middleware.Protected(), uploadHandler.UploadImage)

	// --- Auth ---
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)

	// Strict Rate Limiter for Login (Anti Brute-Force)
	auth.Post("/login", middleware.AuthRateLimiter(), userHandler.Login)

	// --- Users ---
	users := api.Group("/users", middleware.Protected())
	users.Get("/me", userHandler.GetMyProfile)
	users.Put("/me", userHandler.UpdateProfile)

	// --- Toko ---
	stores := api.Group("/toko")
	storesProtected := stores.Group("/", middleware.Protected())
	
	// Define /me BEFORE /:id so it doesn't get trapped by the wildcard
	storesProtected.Get("/me", tokoHandler.GetMyToko)
	storesProtected.Put("/me", tokoHandler.UpdateMyToko)

	stores.Get("/:id", tokoHandler.GetTokoByID)

	// --- Alamat ---
	alamats := api.Group("/alamat", middleware.Protected())
	alamats.Post("/", alamatHandler.CreateAlamat)
	alamats.Get("/", alamatHandler.GetMyAlamat)
	alamats.Get("/:id", alamatHandler.GetAlamatByID)
	alamats.Put("/:id", alamatHandler.UpdateAlamat)
	alamats.Delete("/:id", alamatHandler.DeleteAlamat)

	// --- Category ---
	categories := api.Group("/categories")
	categories.Get("/", categoryHandler.GetAllCategories)
	categories.Get("/:id", categoryHandler.GetCategoryByID)

	categoriesProtected := categories.Group("/", middleware.Protected(), middleware.AdminOnly())
	categoriesProtected.Post("/", categoryHandler.CreateCategory)
	categoriesProtected.Put("/:id", categoryHandler.UpdateCategory)
	categoriesProtected.Delete("/:id", categoryHandler.DeleteCategory)

	// --- Produk ---
	products := api.Group("/produk")
	products.Get("/", produkHandler.GetAllProduk)
	products.Get("/:id", produkHandler.GetProdukByID)

	productsProtected := products.Group("/", middleware.Protected())
	productsProtected.Post("/", produkHandler.CreateProduk)
	productsProtected.Put("/:id", produkHandler.UpdateProduk)

	// --- Transaction ---
	transactions := api.Group("/transactions", middleware.Protected())
	transactions.Post("/checkout", transactionHandler.Checkout)
	transactions.Get("/", transactionHandler.GetMyTransactions)
	transactions.Get("/:id", transactionHandler.GetTransactionDetail)

	return app
}