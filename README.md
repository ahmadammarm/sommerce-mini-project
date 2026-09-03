# Sommerce Mini Project

Sommerce Mini Project is a comprehensive RESTful API built for a multi-vendor e-commerce platform. It provides a robust backend system designed to handle user authentication, vendor store management, product catalogs, and transactional workflows.

This project was developed as a final exam for the Project-Based Internship Program by Evermos, strictly adhering to Clean Architecture principles and modern software engineering standards.

## Tech Stack

*   **Language:** Go (Golang)
*   **Web Framework:** Fiber
*   **ORM:** GORM
*   **Database:** MySQL
*   **Authentication:** JSON Web Tokens (JWT)

## Core Features

### 1. Authentication & Authorization
*   Secure user registration and login workflows using JWT.
*   Unique constraints enforced on user emails and phone numbers.
*   Strict Data Isolation: Users can only retrieve, update, or delete their own data (addresses, store details, products, and transactions).
*   Role-Based Access Control (RBAC): Category management is restricted exclusively to Administrators.

### 2. Multi-Vendor Store System
*   Automated Store Creation: A vendor store is automatically generated and linked to a user immediately upon successful registration.
*   Store Profile Management: Vendors can update their store name and upload store banner images.

### 3. Product Catalog
*   Dual Pricing Support: Products support both reseller pricing and standard consumer pricing.
*   Stock Management: Inventory is automatically validated and deducted during checkout.
*   Rich Media: Supports file uploads for linking multiple images to a single product.
*   Advanced Listings: Implements pagination and filtering on product retrieval endpoints.

### 4. Transaction & Checkout Workflow
*   Shipping Integration: Transactions require a linked user shipping address.
*   Historical Snapshotting: During checkout, the purchased product's current state is snapshotted into a `log_produk` table. This ensures that the transaction receipt remains immutable, even if the vendor alters the live product details in the future.

## Software Architecture

This project implements **Clean Architecture** to separate concerns and ensure maintainability. The codebase is organized into distinct layers:
*   **Delivery (Controllers):** Handles HTTP requests, parses payloads, and returns structured JSON responses.
*   **Use Cases (Services):** Contains the core business logic and rules.
*   **Repositories:** Manages all direct interactions with the MySQL database.

## Documentation

For a detailed breakdown of the database schema and table relationships, please refer to the documentation:
*   [Database Design Documentation](docs/database_design.md)

## Getting Started

### Prerequisites
*   Go 1.24 or higher
*   MySQL Server

### Installation

1. Clone the repository:
```bash
git clone https://github.com/ahmadammarm/sommerce-mini-project.git
cd sommerce-mini-project
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure Environment Variables:
Create a `.env` file in the root directory and configure your database and JWT secrets:
```env
DB_USER=root
DB_PASSWORD=yourpassword
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=sommerce
JWT_SECRET=your_super_secret_key
```

4. Run the application:
```bash
go run cmd/main.go
```
