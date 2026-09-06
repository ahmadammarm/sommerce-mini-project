package alamat

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// AlamatRepository defines the contract for Alamat database operations
type AlamatRepository interface {
	Create(ctx context.Context, alamat *Alamat) error
	FindByUserID(ctx context.Context, userID string) ([]Alamat, error)
	FindByID(ctx context.Context, id string) (*Alamat, error)
	Update(ctx context.Context, alamat *Alamat) error
	Delete(ctx context.Context, id string) error
}

type alamatRepository struct {
	db *gorm.DB
}

// NewAlamatRepository is the constructor for Dependency Injection
func NewAlamatRepository(db *gorm.DB) AlamatRepository {
	return &alamatRepository{db: db}
}

// Create inserts a new address record into the database
func (r *alamatRepository) Create(ctx context.Context, alamat *Alamat) error {
	return r.db.WithContext(ctx).Create(alamat).Error
}

// FindByUserID retrieves all addresses belonging to a specific user
func (r *alamatRepository) FindByUserID(ctx context.Context, userID string) ([]Alamat, error) {
	var alamats []Alamat
	err := r.db.WithContext(ctx).Where("id_user = ?", userID).Find(&alamats).Error
	if err != nil {
		return nil, err
	}
	return alamats, nil
}

// FindByID retrieves a single address by its CUID
func (r *alamatRepository) FindByID(ctx context.Context, id string) (*Alamat, error) {
	var a Alamat
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Graceful error handling
		}
		return nil, err
	}
	return &a, nil
}

// Update saves changes to an existing address
func (r *alamatRepository) Update(ctx context.Context, alamat *Alamat) error {
	return r.db.WithContext(ctx).Save(alamat).Error
}

// Delete completely removes an address from the database
func (r *alamatRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Alamat{}).Error
}
