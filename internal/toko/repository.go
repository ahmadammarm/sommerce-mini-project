package toko

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// TokoRepository defines the contract for Toko database operations
type TokoRepository interface {
	Create(ctx context.Context, toko *Toko) error
	FindByID(ctx context.Context, id string) (*Toko, error)
	FindByUserID(ctx context.Context, userID string) (*Toko, error)
	Update(ctx context.Context, toko *Toko) error
}

type tokoRepository struct {
	db *gorm.DB
}

// NewTokoRepository is the constructor for Dependency Injection (Google Wire)
func NewTokoRepository(db *gorm.DB) TokoRepository {
	return &tokoRepository{db: db}
}

// Create inserts a new store record into the database.
// This will be called immediately after a user registers.
func (r *tokoRepository) Create(ctx context.Context, toko *Toko) error {
	return r.db.WithContext(ctx).Create(toko).Error
}

// FindByID retrieves a store by its unique CUID.
func (r *tokoRepository) FindByID(ctx context.Context, id string) (*Toko, error) {
	var t Toko
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Handled gracefully, returns nil instead of erroring out
		}
		return nil, err
	}
	return &t, nil
}

// FindByUserID retrieves the store that belongs to a specific user.
// This is crucial for verifying ownership before updates.
func (r *tokoRepository) FindByUserID(ctx context.Context, userID string) (*Toko, error) {
	var t Toko
	err := r.db.WithContext(ctx).Where("id_user = ?", userID).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// Update saves changes (like store name or photo URL) to the database.
func (r *tokoRepository) Update(ctx context.Context, toko *Toko) error {
	return r.db.WithContext(ctx).Save(toko).Error
}
