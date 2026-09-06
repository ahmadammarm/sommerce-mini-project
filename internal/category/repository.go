package category

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// CategoryRepository defines the contract for Category database operations
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	FindAll(ctx context.Context) ([]Category, error)
	FindByID(ctx context.Context, id string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id string) error
}

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository is the constructor for Dependency Injection
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create inserts a new category record into the database
func (r *categoryRepository) Create(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// FindAll retrieves every category in the database
func (r *categoryRepository) FindAll(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByID retrieves a single category by its CUID
func (r *categoryRepository) FindByID(ctx context.Context, id string) (*Category, error) {
	var c Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Graceful error handling
		}
		return nil, err
	}
	return &c, nil
}

// Update saves changes to an existing category
func (r *categoryRepository) Update(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete completely removes a category from the database
func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Category{}).Error
}
