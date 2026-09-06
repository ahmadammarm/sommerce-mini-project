package transaction

import (
	"context"

	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
	"gorm.io/gorm"
)

// TransactionRepository defines the contract for Checkout and Transaction database operations
type TransactionRepository interface {
	CreateCheckout(ctx context.Context, trx *Trx, details []DetailTrx, productsToUpdate []produk.Produk) error
	FindByUserID(ctx context.Context, userID string) ([]Trx, error)
	FindDetailsByTrxID(ctx context.Context, trxID string) ([]DetailTrx, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository is the constructor for Dependency Injection
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// CreateCheckout guarantees that Order headers, Details, and Stock decrements happen atomically
func (r *transactionRepository) CreateCheckout(ctx context.Context, trx *Trx, details []DetailTrx, productsToUpdate []produk.Produk) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Insert Header
		if err := tx.Create(trx).Error; err != nil {
			return err // Triggers rollback
		}

		// 2. Insert Transaction Details
		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return err // Triggers rollback
			}
		}

		// 3. Update Product Stocks
		for _, p := range productsToUpdate {
			if err := tx.Save(&p).Error; err != nil {
				return err // Triggers rollback if stock update fails
			}
		}

		return nil // Auto Commits successfully
	})
}

// FindByUserID retrieves all transaction headers for a specific user
func (r *transactionRepository) FindByUserID(ctx context.Context, userID string) ([]Trx, error) {
	var trxs []Trx
	err := r.db.WithContext(ctx).Where("id_user = ?", userID).Find(&trxs).Error
	if err != nil {
		return nil, err
	}
	return trxs, nil
}

// FindDetailsByTrxID retrieves all the items/details purchased within a specific transaction
func (r *transactionRepository) FindDetailsByTrxID(ctx context.Context, trxID string) ([]DetailTrx, error) {
	var details []DetailTrx
	err := r.db.WithContext(ctx).Where("id_trx = ?", trxID).Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}
