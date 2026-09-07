# Sommerce (Social Commerce) - Backend API

Sommerce (kependekan dari Social Commerce) adalah sebuah RESTful API E-Commerce yang dikembangkan sebagai tugas akhir atau proyek mini di program project based internship dari Evermos. 

Sistem ini dirancang dengan standar level produksi menggunakan Clean Architecture (Handler, Service, Repository), menjamin skalabilitas, keamanan, dan keandalan data transaksi finansial.

---

## Fitur Utama

Aplikasi ini telah memenuhi seluruh 15 spesifikasi rubrik yang ketat:
1. Otentikasi & Isolasi Data (JWT): Keamanan identitas mutlak. Pengguna hanya dapat mengakses, mengubah, dan bertransaksi menggunakan data mereka sendiri. Data bersifat unique (Email & No. Telepon).
2. Otomatisasi Entitas: Setiap pengguna yang mendaftar akan otomatis dibuatkan sebuah Toko kosong yang siap digunakan.
3. Manajemen Kategori Terpusat: Hanya peran Admin yang memiliki wewenang untuk menambah, mengubah, atau menghapus kategori produk.
4. API Upload File Murni: Mendukung pengunggahan gambar (multipart/form-data) dengan validasi ukuran (Maks 5MB) dan ekstensi file (JPG/PNG).
5. Validasi Alamat Integrasi Emsifa: Alamat pengiriman pengguna divalidasi secara real-time ke API Emsifa v2.
6. Checkout Transaksi Anti-Bocor (ACID): Menggunakan Pessimistic Locking dan Idempotency Key di database MySQL untuk menjamin mustahilnya terjadi stok minus atau checkout ganda.
7. Snapshot Log Produk Historis: Setiap kali transaksi terjadi, sistem akan mencetak log produk statis. Ini memastikan kuitansi transaksi selalu mengikat pada harga dan nama produk di detik pembelian, meskipun harga produk diubah di masa depan.
8. Pagination & Filtering: Pengambilan data produk dilengkapi dengan sistem limitasi halaman dan pencarian.

---

## Teknologi yang Digunakan

* Bahasa: Golang (Go 1.24+)
* Framework Web: Go Fiber v2
* Database & ORM: MySQL & GORM
* Dependency Injection: Google Wire
* Utilitas: CUID (Unique ID), Golang-JWT, Joho/Godotenv

---

## Struktur Direktori

Proyek ini sangat mematuhi prinsip Clean Architecture. Setiap domain bisnis dipisah secara mandiri dan tidak saling tumpang tindih.

```text
.
|-- cmd/                    # Titik masuk (Entrypoint) utama aplikasi
|   |-- api/                # File inisialisasi server (main.go) dan injeksi dependensi (wire.go)
|   +-- seeder/             # Skrip eksekusi mandiri untuk mengisi database awal
|-- docs/                   # Direktori kumpulan dokumentasi sistem dan database
|   |-- alur_sistem.md      # Panduan alur registrasi, login, hingga transaksi
|   |-- database_design.md  # Penjelasan skema tabel dan sistem penguncian pesimis
|   +-- sommerce_collection.json # Berkas JSON API Postman siap impor
|-- internal/               # Logika bisnis inti yang terisolasi berdasarkan domain
|   |-- alamat/             # Domain manajemen alamat pengiriman
|   |-- category/           # Domain manajemen kategori (khusus Admin)
|   |-- middleware/         # Penengah request (Verifikasi JWT, Security, Rate Limiter)
|   |-- produk/             # Domain manajemen produk dan foto
|   |-- router/             # Titik kumpul pendaftaran seluruh rute (Endpoint API)
|   |-- toko/               # Domain profil toko bawaan
|   |-- transaction/        # Domain transaksi (Checkout, Pessimistic Locking, Detail)
|   |-- upload/             # Domain khusus penerimaan file unggahan
|   +-- user/               # Domain registrasi, login, dan profil pengguna
|-- pkg/                    # Pustaka utilitas (Reusable packages) untuk berbagai domain
|   |-- config/             # Konfigurasi koneksi Database MySQL dan Validator struct
|   |-- emsifa/             # Klien HTTP kustom untuk mengecek validitas ID ke API Emsifa
|   |-- hash/               # Enkripsi dan dekripsi kata sandi menggunakan algoritma Bcrypt
|   |-- response/           # Format pembakuan struktur JSON balasan (WebResponse)
|   |-- uid/                # Generator pengenal unik (CUID)
|   +-- upload/             # Fungsi validasi ekstensi dan ukuran file gambar
|-- public/                 # Folder publik yang bisa diakses langsung via web
|   +-- uploads/            # Ruang penyimpanan fisik untuk gambar yang diunggah
|-- utils/                  # Fungsi bantuan operasional (seperti pembuatan Token JWT)
|-- .env                    # Berkas penyimpan variabel rahasia (Kredensial DB & JWT Secret)
+-- README.md               # Dokumentasi utama proyek
```

Catatan Arsitektur Internal:
Di dalam setiap folder domain (misalnya `internal/user/`), lapisan dipisahkan kembali menjadi:
* `entity.go`: Cetak biru struktur database (GORM).
* `dto.go`: Struktur data serah-terima antara klien dan server (JSON).
* `repository.go`: Lapisan yang berinteraksi langsung dengan MySQL.
* `service.go`: Lapisan yang memegang seluruh otak dan logika bisnis.
* `handler.go`: Lapisan terluar yang menerjemahkan Request Fiber menjadi fungsi Service.
* `wire.go`: Lapisan penyedia injeksi untuk Google Wire.

---

## Persyaratan & Instalasi

### 1. Persiapan Database
Pastikan Anda memiliki MySQL yang berjalan (misal: menggunakan Laragon). Buat sebuah database kosong, contohnya bernama `sommerce`.

### 2. Konfigurasi Environment (.env)
Buat file bernama `.env` di root folder proyek, dan isi dengan variabel berikut:
```env
PORT=3000
DATABASE_URL=root:@tcp(127.0.0.1:3306)/sommerce?charset=utf8mb4&parseTime=True&loc=Local
JWT_SECRET=super-secret-sommerce-key
```
(Sesuaikan password root dan nama database jika berbeda).

### 3. Menjalankan Aplikasi
Buka terminal dan jalankan perintah:
```bash
go mod tidy
go run ./cmd/api
```
(Sistem akan otomatis melakukan Migrasi Tabel Database / AutoMigrate saat pertama kali menyala).

### 4. Menjalankan Database Seeder (Opsional tapi Disarankan)
Untuk mengisi database dengan akun Admin, Toko Admin, Kategori awal, dan Produk bawaan, buka terminal baru dan jalankan:
```bash
go run ./cmd/seeder
```
* Email Admin: admin@sommerce.com
* Password Admin: admin123

---

## Dokumentasi Lanjutan (Direktori docs)

Untuk membedah sistem lebih dalam, kami telah menyediakan kumpulan dokumen teknis terpusat di dalam direktori `docs/`. Silakan merujuk pada:

1. **Struktur Database (docs/database_design.md):** Penjelasan skema tabel, relasi entitas, dan logika perlindungan stok transaksi.
2. **Alur Sistem (docs/alur_sistem.md):** Penjelasan logika bisnis dan aliran data tiap modul.
3. **Dokumentasi API (docs/sommerce_collection.json):** Berkas dokumentasi endpoints utuh.

**Cara Membaca Dokumentasi API:**
1. Buka aplikasi Postman.
2. Klik tombol Import (atau Ctrl+O).
3. Pilih dan seret file `sommerce_collection.json` dari dalam folder `docs/`.
4. Sebuah koleksi bernama "Sommerce API" akan muncul.
5. Penting: Rute di dalam koleksi ini dilengkapi dengan Script Automasi. Cukup jalankan endpoint Login, dan token Anda akan otomatis terpasang ke seluruh endpoint lainnya tanpa perlu copy-paste!

### Catatan Khusus: Domain Alamat (Emsifa API)
Pembuatan Alamat Pengiriman (POST /api/alamat) sengaja dirancang secara ketat (strict) untuk hanya menerima ID Provinsi dan ID Kota sesuai dengan standar API Wilayah Indonesia (Emsifa). 
* Contoh Provinsi Jawa Timur: "35"
* Contoh Kota Malang: "35.73"

Ini dipertahankan untuk memastikan integritas data spasial tetap akurat dan relevan dengan sistem logistik sungguhan. Sistem backend akan menembak ke URL Emsifa untuk memverifikasi apakah ID yang dikirim oleh pengguna benar-benar ada di Indonesia. Jika ID tidak valid, pembuatan alamat akan otomatis ditolak.

---

## Keamanan & Batasan (Rate Limiting)
Aplikasi ini dilindungi oleh Global Rate Limiter untuk menangkis serangan DDoS, dan Strict Auth Limiter khusus pada rute `/auth/login` untuk mencegah serangan tebak password (Brute-force). Seluruh rute otentikasi mensyaratkan format header Authorization: Bearer <token>.
