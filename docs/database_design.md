# Dokumentasi Desain Database Sommerce

Dokumen ini menjabarkan skema database untuk Proyek Mini Sommerce. Sistem ini dirancang sebagai platform e-commerce yang mendukung kapabilitas multi-vendor serta memisahkan harga untuk reseller dan konsumen umum.

## Gambaran Umum
Skema database ini mendukung sebuah sistem di mana pengguna dapat berperan ganda sebagai pembeli maupun penjual (vendor). Sistem ini dirancang dengan penanganan kokoh terhadap variasi harga produk, pemenuhan keranjang belanja dari berbagai vendor secara bersamaan, serta riwayat transaksi yang bersifat *immutable* melalui mekanisme pencatatan log produk.

---

## 1. Manajemen Pengguna

### Tabel `users`
Tabel sentral untuk manajemen akun dan otentikasi.
*   **Primary Key**: `id` (CUID)
*   **Fields**:
    *   `nama`, `email`, `kata_sandi`: Profil dasar dan kredensial akses masuk. Kolom `email` bersifat `UNIQUE`.
    *   `notelp`: Nomor telepon dengan batasan `UNIQUE` constraint untuk mencegah pendaftaran ganda.
    *   `tanggal_lahir`, `jenis_kelamin`, `tentang`, `pekerjaan`: Detail profil tambahan milik pengguna.
    *   `id_provinsi`, `id_kota`: Referensi ke data statis wilayah Indonesia (melalui Emsifa API) untuk menentukan domisili utama pengguna.
    *   `isAdmin`: Boolean flag untuk membedakan antara administrator sistem dan pengguna biasa.
    *   `created_at`, `updated_at`: Pencatat waktu otomatis.

### Tabel `alamats` (Alamat Pengiriman)
Menyimpan data alamat fisik yang terikat dengan pengguna untuk keperluan pengiriman barang.
*   **Primary Key**: `id`
*   **Foreign Key**: `id_user` menunjuk ke `users(id)`
*   **Fields**: `judul_alamat` (contoh: Rumah, Kantor), `nama_penerima`, `no_telp`, `detail_alamat`, `id_provinsi`, `id_kota`.
*   **Relasi**: **1-to-Many**. Satu pengguna (`users`) dapat menyimpan dan mengelola banyak alamat (`alamats`) untuk dipilih pada saat proses checkout.

---

## 2. Sistem Toko dan Penjual

### Tabel `tokos`
Merepresentasikan profil etalase toko yang dimiliki oleh seorang pengguna di dalam platform.
*   **Primary Key**: `id`
*   **Foreign Key**: `id_user` menunjuk ke `users(id)` (`UNIQUE`)
*   **Fields**: `nama_toko`, `url_foto` (URL untuk logo atau banner toko).
*   **Relasi**: **1-to-1**. Setiap akun pengguna hanya dapat memiliki tepat satu toko. Entitas toko ini dibuat secara otomatis oleh sistem saat pengguna pertama kali mendaftar.

---

## 3. Katalog Produk

### Tabel `categories`
Tabel *lookup* untuk pengelompokan jenis produk. Dikelola sepenuhnya oleh Administrator.
*   **Primary Key**: `id`
*   **Fields**: `nama_category` (Nama Kategori)

### Tabel `produks`
Menyimpan daftar produk aktif yang tersedia untuk dibeli.
*   **Primary Key**: `id`
*   **Foreign Keys**: 
    *   `id_toko` menunjuk ke `tokos(id)`
    *   `id_category` menunjuk ke `categories(id)`
*   **Fields**: `nama_produk`, `slug` (teks SEO-friendly untuk URL), `harga_reseller` (harga khusus pengecer), `harga_konsumen` (harga ritel standar), `stok` (jumlah ketersediaan barang), `deskripsi`.
*   **Relasi**: 
    *   **Many-to-1** dengan `tokos`: Satu toko dapat mendaftarkan banyak produk.
    *   **Many-to-1** dengan `categories`: Satu kategori dapat menampung banyak produk.

### Tabel `foto_produks`
Menyimpan galeri gambar untuk setiap produk.
*   **Primary Key**: `id`
*   **Foreign Key**: `id_produk` menunjuk ke `produks(id)`
*   **Fields**: `url` (Lokasi publik dari berkas gambar yang telah diunggah).
*   **Relasi**: **1-to-Many**. Satu produk (`produks`) dapat memiliki banyak foto (`foto_produks`).

### Tabel `log_produks` (Rekam Jejak Historis Produk)
Mempertahankan *snapshot* historis dari rincian produk pada detik tertentu.
*   **Primary Key**: `id`
*   **Fields**: Merupakan cerminan pasti dari tabel `produks` (`id_produk`, `nama_produk`, `slug`, `harga_reseller`, `harga_konsumen`, `deskripsi`, `id_toko`, `id_category`).
*   **Tujuan**: *Snapshot* historis ini dibuat dan diisi secara eksklusif hanya pada saat transaksi terjadi (saat *checkout*). Bukti kuitansi transaksi diikat ke tabel ini, sehingga riwayat transaksi terjamin integritasnya.

---

## 4. Pesanan dan Transaksi

### Tabel `trxs` (Header Transaksi)
Merekam satu sesi penuh pemesanan atau keranjang belanja dari seorang pembeli.
*   **Primary Key**: `id`
*   **Foreign Keys**:
    *   `id_user` menunjuk ke `users(id)`
    *   `alamat_pengiriman` menunjuk ke `alamats(id)`
*   **Fields**: `harga_total` (total keseluruhan yang harus dibayar), `kode_invoice` (pengenal kuitansi unik), `method_bayar` (metode pembayaran).
*   **Idempotency Key**: Kolom dengan `UNIQUE constraint` yang mencegah pengguna secara tidak sengaja memproses keranjang yang sama dua kali akibat gangguan jaringan (Double Checkout).
*   **Relasi**: **1-to-Many** (`users` memiliki banyak `trxs`).

### Tabel `detail_trxs` (Rincian Item Transaksi)
Merekam barang-barang individual yang dibeli di dalam suatu transaksi.
*   **Primary Key**: `id`
*   **Foreign Keys**:
    *   `id_trx` menunjuk ke `trxs(id)`
    *   `id_log_produk` menunjuk ke `log_produks(id)`
    *   `id_toko` menunjuk ke `tokos(id)`
*   **Fields**: `kuantitas` (jumlah yang dibeli), `harga_total` (subtotal untuk rincian ini).
*   **Relasi**:
    *   **Many-to-1** dengan `trxs`: Satu pesanan memiliki banyak rincian barang.
    *   **Many-to-1** dengan `log_produks`: Mengikat rincian pesanan dengan *snapshot* historis, memastikan nota belanja tidak akan berubah meskipun harga produk asli diubah oleh penjual di masa depan.
    *   **Many-to-1** dengan `tokos`: Melacak toko mana yang harus memenuhi (mengirimkan) rincian pesanan spesifik ini.

---

## Keputusan Arsitektur Kunci

1.  **Rekam Jejak Finansial Permanen (Immutable)**: Dengan menautkan rincian transaksi (`detail_trxs`) ke tabel *snapshot* (`log_produks`) alih-alih ke tabel produk yang aktif, platform memastikan bahwa perubahan harga atau deskripsi yang dilakukan oleh vendor tidak akan pernah merusak akurasi data pesanan di masa lampau.
2.  **Perlindungan Konkurensi (Pessimistic Locking)**: Pada saat transaksi dibuat, sistem menerapkan perintah `SELECT ... FOR UPDATE` pada tabel `produks`. Ini memaksa pembeli lain untuk mengantre di tingkat database (menerapkan *Row Lock*), menjamin bahwa stok barang tidak akan pernah bernilai negatif (minus) walau diakses ribuan pengguna secara serentak (*Race Condition*).
3.  **Dukungan Pengecer (Reseller)**: Penyertaan atribut ganda `harga_reseller` dan `harga_konsumen` secara langsung mendukung model bisnis waralaba, di mana reseller terdaftar berhak mendapatkan potongan harga sistemik dibandingkan dengan pengguna standar.
4.  **Satu Keranjang Berbagai Toko (Multi-Vendor Cart)**: Mengingat tabel `detail_trxs` menyimpan nilai `id_toko` secara mandiri pada setiap baris item, struktur database ini secara alami mendukung skenario *multi-vendor* di mana pengguna membeli banyak barang dari berbagai toko yang berbeda di dalam satu kali proses checkout.
5.  **Optimalisasi Indeks (Database Indexing)**: Kinerja *query* diperkuat secara signifikan dengan penerapan `index` dan `uniqueIndex` (B-Tree) pada mayoritas *Foreign Key* dan kolom pencarian. Hal ini menjamin bahwa operasi *JOIN*, filter produk (`id_category`, `id_toko`), pencarian nama (`nama_produk`), serta validasi otorisasi kepemilikan keranjang tetap beroperasi secara *O(log N)* alih-alih *Full Table Scan* yang memberatkan peladen.
