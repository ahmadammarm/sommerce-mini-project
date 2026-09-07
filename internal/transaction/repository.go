package transaction

import (
	"context"
	"fmt"

	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/uid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	CreateCheckout(ctx context.Context, trx *Trx, items []CheckoutItemRequest) error
	FindByUserID(ctx context.Context, userID string) ([]Trx, error)
	FindDetailsByTrxID(ctx context.Context, trxID string) ([]DetailTrx, error)
}

type transactionRepository struct {
	db    *gorm.DB
	idGen uid.IDGenerator
}

func NewTransactionRepository(db *gorm.DB, idg uid.IDGenerator) TransactionRepository {
	return &transactionRepository{db: db, idGen: idg}
}

func (r *transactionRepository) CreateCheckout(ctx context.Context, trx *Trx, items []CheckoutItemRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var details []DetailTrx
		var grandTotal int

		for _, item := range items {
			var p produk.Produk
			
			// PESSIMISTIC LOCK: Lock the product row so no other checkout can read it simultaneously
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", item.IdProduk).First(&p).Error; err != nil {
				return fmt.Errorf("product %s not found or locked", item.IdProduk)
			}

			// Validate Stock
			if p.Stok < item.Kuantitas {
				return fmt.Errorf("insufficient stock for product %s (available: %d)", p.NamaProduk, p.Stok)
			}

			// Deduct Stock
			p.Stok -= item.Kuantitas
			if err := tx.Save(&p).Error; err != nil {
				return err
			}

			subTotal := p.HargaKonsumen * item.Kuantitas
			grandTotal += subTotal

			details = append(details, DetailTrx{
				ID:          r.idGen.GenerateID(),
				IdTrx:       trx.ID,
				IdLogProduk: p.ID, 
				IdToko:      p.IdToko,
				Kuantitas:   item.Kuantitas,
				HargaTotal:  subTotal,
			})
		}

		trx.HargaTotal = grandTotal

		// Save Header
		if err := tx.Create(trx).Error; err != nil {
			return err
		}

		// Save Details
		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *transactionRepository) FindByUserID(ctx context.Context, userID string) ([]Trx, error) {
	var trxs []Trx
	err := r.db.WithContext(ctx).Where("id_user = ?", userID).Find(&trxs).Error
	return trxs, err
}

func (r *transactionRepository) FindDetailsByTrxID(ctx context.Context, trxID string) ([]DetailTrx, error) {
	var details []DetailTrx
	err := r.db.WithContext(ctx).Where("id_trx = ?", trxID).Find(&details).Error
	return details, err
}
