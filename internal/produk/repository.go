package produk

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// ProdukRepository defines the contract for Produk database operations
type ProdukRepository interface {
	Create(ctx context.Context, produk *Produk, fotos []FotoProduk) error
	Update(ctx context.Context, produk *Produk) error
	FindByID(ctx context.Context, id string) (*Produk, error)
	FindFotosByProdukID(ctx context.Context, produkID string) ([]FotoProduk, error)
	FindAll(ctx context.Context, filter ProductFilterDTO) ([]Produk, int64, error)
}

type produkRepository struct {
	db *gorm.DB
}

// NewProdukRepository is the constructor for Dependency Injection
func NewProdukRepository(db *gorm.DB) ProdukRepository {
	return &produkRepository{db: db}
}

// Create uses a DB transaction to save the Product and Photos atomically
func (r *produkRepository) Create(ctx context.Context, produk *Produk, fotos []FotoProduk) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Insert Main Product
		if err := tx.Create(produk).Error; err != nil {
			return err // Triggers rollback
		}

		// 2. Insert Photos (if any)
		if len(fotos) > 0 {
			if err := tx.Create(&fotos).Error; err != nil {
				return err
			}
		}

		return nil // Commits successfully
	})
}

// Update saves the Product directly
func (r *produkRepository) Update(ctx context.Context, produk *Produk) error {
	return r.db.WithContext(ctx).Save(produk).Error
}

// FindByID retrieves a single product by its CUID
func (r *produkRepository) FindByID(ctx context.Context, id string) (*Produk, error) {
	var p Produk
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Graceful error handling
		}
		return nil, err
	}
	return &p, nil
}

// FindFotosByProdukID retrieves all photos associated with a product
func (r *produkRepository) FindFotosByProdukID(ctx context.Context, produkID string) ([]FotoProduk, error) {
	var fotos []FotoProduk
	err := r.db.WithContext(ctx).Where("id_produk = ?", produkID).Find(&fotos).Error
	if err != nil {
		return nil, err
	}
	return fotos, nil
}

// FindAll dynamically builds the SQL query based on Search, Category filtering, and Pagination
func (r *produkRepository) FindAll(ctx context.Context, filter ProductFilterDTO) ([]Produk, int64, error) {
	var produks []Produk
	var totalItems int64

	// Start building the query
	query := r.db.WithContext(ctx).Model(&Produk{})

	// Apply Dynamic Filters
	if filter.Search != "" {
		query = query.Where("nama_produk LIKE ?", "%"+filter.Search+"%")
	}
	if filter.CategoryId != "" {
		query = query.Where("id_category = ?", filter.CategoryId)
	}

	// Count total matching items before applying pagination offsets
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Apply Pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Execute final retrieval query
	if err := query.Find(&produks).Error; err != nil {
		return nil, 0, err
	}

	return produks, totalItems, nil
}
