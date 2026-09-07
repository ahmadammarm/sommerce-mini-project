package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ahmadammarm/sommerce-mini-project/internal/category"
	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
	"github.com/ahmadammarm/sommerce-mini-project/internal/toko"
	"github.com/ahmadammarm/sommerce-mini-project/internal/user"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/config"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/hash"
	"github.com/joho/godotenv"
	"github.com/lucsky/cuid"
)

func main() {
	// 1. Load env dan connect DB
	godotenv.Load()
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatalf("Gagal konek DB: %v", err)
	}

	fmt.Println("🚀 Memulai proses Seeding Database...")

	// Cek apakah admin sudah ada agar tidak error duplikat saat dijalankan berulang kali
	var count int64
	db.Model(&user.User{}).Where("email = ?", "admin@sommerce.com").Count(&count)
	if count > 0 {
		fmt.Println("✅ Data sudah di-seed sebelumnya. Menghentikan proses seeding.")
		return
	}

	// 2. SEED ADMIN USER
	adminID := cuid.New()
	hashedPassword, _ := hash.HashPassword("admin123")
	adminDate, _ := time.Parse("2006-01-02", "1990-01-01")

	adminUser := user.User{
		ID:           adminID,
		Nama:         "Super Admin Sommerce",
		Email:        "admin@sommerce.com",
		KataSandi:    hashedPassword,
		NoTelp:       "080000000000",
		TanggalLahir: adminDate,
		JenisKelamin: "Laki-laki",
		Tentang:      "Saya adalah penguasa sistem ini",
		Pekerjaan:    "System Administrator",
		IdProvinsi:   "11",   // Aceh
		IdKota:       "1101", // Kab. Simeulue
		IsAdmin:      true,   // Kunci akses ke endpoint Category
	}
	db.Create(&adminUser)
	fmt.Println("👤 Admin User berhasil dibuat.")

	// 3. SEED TOKO ADMIN
	tokoID := cuid.New()
	adminToko := toko.Toko{
		ID:       tokoID,
		IdUser:   adminID,
		NamaToko: "Toko Sommerce Official",
		UrlFoto:  "https://example.com/logo-sommerce.png",
	}
	db.Create(&adminToko)
	fmt.Println("🏪 Toko Admin berhasil dibuat.")

	// 4. SEED KATEGORI
	catElektronikID := cuid.New()
	catPakaianID := cuid.New()
	catMakananID := cuid.New()

	db.Create([]category.Category{
		{ID: catElektronikID, NamaCategory: "Elektronik"},
		{ID: catPakaianID, NamaCategory: "Fashion & Pakaian"},
		{ID: catMakananID, NamaCategory: "Makanan & Minuman"},
	})
	fmt.Println("🗂️ 3 Kategori berhasil dibuat.")

	// 5. SEED PRODUK
	produkID := cuid.New()
	p := produk.Produk{
		ID:            produkID,
		IdToko:        tokoID,
		IdCategory:    catElektronikID,
		NamaProduk:    "Laptop Pro M3 15-inch",
		Slug:          "laptop-pro-m3-15-inch",
		HargaReseller: 18000000,
		HargaKonsumen: 20000000,
		Stok:          50,
		Deskripsi:     "Laptop super canggih untuk memprogram backend Go.",
	}
	db.Create(&p)
	fmt.Println("💻 Produk berhasil dibuat.")

	// 6. SEED LOG PRODUK (Mewujudkan Rule 8)
	logP := produk.LogProduk{
		ID:            cuid.New(),
		IdProduk:      produkID,
		IdToko:        tokoID,
		IdCategory:    catElektronikID,
		NamaProduk:    p.NamaProduk,
		Slug:          p.Slug,
		HargaReseller: p.HargaReseller,
		HargaKonsumen: p.HargaKonsumen,
		Deskripsi:     p.Deskripsi,
	}
	db.Create(&logP)
	fmt.Println("📸 Snapshot Log Produk berhasil disimpan.")

	fmt.Println("\n🎉 SEEDING SELESAI!")
	fmt.Println("=====================================")
	fmt.Println("Gunakan kredensial ini untuk login:")
	fmt.Println("Email    : admin@sommerce.com")
	fmt.Println("Password : admin123")
	fmt.Println("=====================================")
}
