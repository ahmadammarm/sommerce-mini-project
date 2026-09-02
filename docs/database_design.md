# Sommerce Mini Project - Database Design Documentation

This document outlines the database schema for the Sommerce Mini Project, which is designed as an e-commerce/marketplace platform with multi-vendor and reseller capabilities.

## Overview
The schema supports a system where users can act as buyers or sellers (vendors). It includes robust handling for product variations in pricing (reseller vs. consumer), multi-vendor cart fulfillment, and immutable transaction history through product logging.

---

## 1. User Management

### `User` Table
The central table for account management and authentication.
*   **Primary Key**: `id`
*   **Fields**:
    *   `nama`, `email`, `kata_sandi` (password): Basic profile and login credentials.
    *   `notelp`: Unique phone number constraint.
    *   `tanggal_lahir` (birth date), `jenis_kelamin` (gender), `tentang` (about), `pekerjaan` (occupation): User profile details.
    *   `id_provinsi`, `id_kota`: References to external/static location tables (not pictured) for the user's primary region.
    *   `isAdmin`: Boolean flag to identify system administrators versus regular users.
    *   Timestamps: `created_at`, `updated_at`

### `alamat` (Address) Table
Stores physical addresses associated with a user for shipping purposes.
*   **Primary Key**: `id`
*   **Foreign Key**: `id_user` -> `User(id)`
*   **Fields**: `judul_alamat` (e.g., Home, Office), `nama_penerima` (recipient name), `no_telp` (contact number), `detail_alamat` (full physical address).
*   **Relationship**: **1-to-Many** (`User` has many `alamat`). A user can store multiple shipping addresses to choose from during checkout.

---

## 2. Store & Vendor System

### `toko` (Store) Table
Represents a shop front owned by a user on the platform.
*   **Primary Key**: `id`
*   **Foreign Key**: `id_user` -> `User(id)`
*   **Fields**: `nama_toko` (Store Name), `url_foto` (Store Logo/Banner URL).
*   **Relationship**: **1-to-1** (`User` has one `toko`). A user account can be associated with a single vendor shop.

---

## 3. Product Catalog

### `category` Table
Lookup table for product categories.
*   **Primary Key**: `id`
*   **Fields**: `nama_category`

### `produk` (Product) Table
Stores active listings available for purchase.
*   **Primary Key**: `id`
*   **Foreign Keys**: 
    *   `id_toko` -> `toko(id)`
    *   `id_category` -> `category(id)`
*   **Fields**: `nama_produk`, `slug` (for SEO-friendly URLs), `harga_reseller` (price for resellers), `harga_konsumen` (standard retail price), `stok` (inventory quantity), `deskripsi`.
*   **Relationships**: 
    *   **Many-to-1** with `toko`: A store lists many products.
    *   **Many-to-1** with `category`: A category contains many products.

### `foto_produk` (Product Photos) Table
Stores image galleries for products.
*   **Primary Key**: `id`
*   **Foreign Key**: `id_produk` -> `produk(id)`
*   **Fields**: `url`
*   **Relationship**: **1-to-Many** (`produk` has many `foto_produk`).

### `log_produk` (Product Log) Table
Maintains historical snapshots of product details.
*   **Primary Key**: `id`
*   **Fields**: Mirrors the `produk` table (`id_produk`, `nama_produk`, `slug`, `harga_reseller`, `harga_konsumen`, `deskripsi`, `id_toko`, `id_category`).
*   **Purpose**: Whenever a product is modified, or when a purchase occurs, a snapshot can be logged here. Transactions link to this table instead of the live `produk` table to preserve historical integrity.

---

## 4. Order & Transactions

### `trx` (Transaction Header) Table
Records an entire checkout session/order placed by a user.
*   **Primary Key**: `id`
*   **Foreign Keys**:
    *   `id_user` -> `User(id)`
    *   `alamat_pengiriman` -> `alamat(id)`
*   **Fields**: `harga_total` (grand total), `kode_invoice` (unique invoice identifier), `method_bayar` (payment method).
*   **Relationship**: **1-to-Many** (`User` has many `trx`).

### `detail_trx` (Transaction Line Items) Table
Records the individual items purchased within a transaction.
*   **Primary Key**: `id`
*   **Foreign Keys**:
    *   `id_trx` -> `trx(id)`
    *   `id_log_produk` -> `log_produk(id)`
    *   `id_toko` -> `toko(id)`
*   **Fields**: `kuantitas` (quantity bought), `harga_total` (subtotal for this line item).
*   **Relationships**:
    *   **Many-to-1** with `trx`: One transaction has many line items.
    *   **Many-to-1** with `log_produk`: Links to the snapshot of the product to ensure the receipt remains accurate even if the live product is altered.
    *   **Many-to-1** with `toko`: Tracks which store fulfills this specific line item.

---

## Key Architectural Decisions

1.  **Immutable Financial Records (`log_produk`)**: By tying transaction details (`detail_trx`) to a product snapshot (`log_produk`) rather than the active `produk` table, the platform ensures that changes to a product's price or description by a vendor will not retroactively alter the data of completed orders.
2.  **Reseller/Dropship Support**: The inclusion of both `harga_reseller` and `harga_konsumen` directly supports a business model where authorized resellers receive discounted rates while standard users pay retail pricing.
3.  **Multi-Vendor Checkout**: Because `detail_trx` tracks `id_toko` at the line-item level, the database natively supports users checking out a single cart (`trx`) containing items from multiple different vendors.
