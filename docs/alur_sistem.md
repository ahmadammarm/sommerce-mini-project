# Alur Sistem Sommerce

Dokumen ini menjelaskan alur logika fungsional utama di dalam aplikasi backend Sommerce, mulai dari pendaftaran pengguna hingga proses transaksi yang kompleks.

## 1. Alur Pendaftaran Pengguna (Registrasi)
Alur ini memastikan bahwa setiap pengguna baru didaftarkan secara aman dan memiliki relasi yang utuh di dalam sistem.

1. Pengguna mengirimkan data registrasi (nama, email, kata sandi, dll).
2. Sistem memvalidasi kelengkapan data menggunakan validator struktur (Go Playground Validator).
3. Sistem memeriksa keunikan **Email** dan **Nomor Telepon** ke database. Jika sudah ada, sistem menolak pendaftaran.
4. Sistem melakukan validasi silang (cross-check) ID Provinsi dan ID Kota yang dikirim pengguna ke **Emsifa API** secara real-time.
5. Sistem mengenkripsi kata sandi menggunakan algoritma **Bcrypt**.
6. Sistem menyimpan data Pengguna (User) ke dalam database.
7. Sistem **secara otomatis** membuatkan entitas **Toko** bawaan untuk pengguna tersebut, menggunakan ID Pengguna sebagai referensi.

## 2. Alur Otentikasi dan Otorisasi (Login)
Alur ini mengamankan rute-rute privat aplikasi menggunakan JWT (JSON Web Token).

1. Pengguna mengirimkan kredensial (Email dan Kata Sandi).
2. Sistem mencari pengguna berdasarkan email. Jika tidak ada, kembalikan pesan error.
3. Sistem membandingkan kata sandi yang di-hash di database dengan masukan pengguna menggunakan Bcrypt.
4. Jika cocok, sistem mencetak Token JWT yang berisi `user_id` dan ditandatangani dengan kunci rahasia rahasia (`JWT_SECRET`).
5. Pada permintaan API selanjutnya, middleware akan mencegat setiap permintaan untuk memverifikasi keaslian token JWT. Jika valid, ID pengguna diekstrak untuk digunakan pada logika selanjutnya (mencegah pengguna memanipulasi data pengguna lain).

## 3. Alur Pembuatan Produk dan Unggah Gambar
Alur ini memisahkan proses pengunggahan berkas fisik dan penyisipan data tekstual.

1. **Tahap 1 (Upload):** Klien menembak endpoint `/api/upload` dengan berkas gambar. Sistem memvalidasi ekstensi (JPG/PNG) dan ukuran (Maksimal 5MB), lalu menyimpannya secara statis dan mengembalikan URL publik.
2. **Tahap 2 (Buat Produk):** Klien menembak endpoint `/api/produk` dengan data detail produk beserta URL gambar (berbentuk array JSON).
3. Sistem memvalidasi bahwa Toko dari pengguna yang sedang login benar-benar ada.
4. Sistem menyimpan data utama Produk ke tabel `produks`.
5. Sistem menyimpan seluruh daftar URL foto ke tabel `foto_produks`.
6. Seluruh proses nomor 4 dan 5 dibungkus dalam *Database Transaction*, sehingga jika gagal di salah satu titik, semuanya akan digagalkan (rollback).

## 4. Alur Transaksi dan Checkout (ACID & Concurrency Safe)
Ini adalah alur paling krusial di aplikasi, dirancang untuk mencegah anomali data saat lalu lintas pengguna sedang tinggi.

1. Pengguna mengirimkan daftar barang yang ingin dibeli, ID Alamat pengiriman, dan **Idempotency Key** di Header.
2. Sistem memvalidasi Idempotency Key. Jika kunci tersebut sudah pernah digunakan untuk transaksi yang sukses, database akan langsung menolak proses tersebut untuk mencegah duplikasi (Double Checkout).
3. Sistem memvalidasi bahwa alamat pengiriman yang dipilih adalah milik pengguna yang sedang login.
4. **Membuka Database Transaction:**
   * Untuk setiap barang di keranjang belanja, sistem menjalankan kueri **SELECT ... FOR UPDATE** (Pessimistic Locking). Ini menerapkan *Row Lock* pada baris produk tersebut di tabel MySQL, memaksa pengguna lain yang ingin membeli produk yang sama untuk menunggu.
   * Sistem memeriksa stok produk saat ini. Jika stok kurang dari kuantitas yang diminta, transaksi langsung digagalkan (Rollback) dan *Row Lock* dilepas.
   * Jika stok cukup, sistem memotong stok secara aman, dan menyimpan perubahannya kembali ke database.
   * Sistem membuat **Snapshot Harga** dengan menyalin data produk saat itu ke tabel `log_produks`.
   * Sistem membuat baris detail transaksi (`detail_trxs`) yang diikat (direlasikan) dengan ID Snapshot Log Produk tadi, bukan ke ID Produk asli.
5. Sistem menjumlahkan semua harga dan menyimpan informasi utama (header) ke tabel `trxs`.
6. Jika tidak ada satu pun kendala sejak langkah 4, sistem akan melakukan **Commit**. Gembok database dilepas dan transaksi resmi tercatat.
