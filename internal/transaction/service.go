package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ahmadammarm/sommerce-mini-project/internal/alamat"
	"github.com/ahmadammarm/sommerce-mini-project/internal/produk"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/uid"
	"github.com/go-playground/validator/v10"
)

type TransactionService interface {
	Checkout(ctx context.Context, userID string, req *CheckoutRequest) (*TransactionResponse, error)
	GetMyTransactions(ctx context.Context, userID string) ([]TransactionResponse, error)
	GetTransactionDetail(ctx context.Context, userID, trxID string) (*TransactionResponse, error) // Can be expanded to return items too
}

type transactionService struct {
	trxRepo    TransactionRepository
	produkRepo produk.ProdukRepository
	alamatRepo alamat.AlamatRepository
	idGen      uid.IDGenerator
	validator  *validator.Validate
}

func NewTransactionService(tr TransactionRepository, pr produk.ProdukRepository, ar alamat.AlamatRepository, idg uid.IDGenerator, v *validator.Validate) TransactionService {
	return &transactionService{
		trxRepo:    tr,
		produkRepo: pr,
		alamatRepo: ar,
		idGen:      idg,
		validator:  v,
	}
}

func (s *transactionService) Checkout(ctx context.Context, userID string, req *CheckoutRequest) (*TransactionResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// 1. Verify Alamat Ownership
	addr, err := s.alamatRepo.FindByID(ctx, req.AlamatPengiriman)
	if err != nil {
		return nil, err
	}
	if addr == nil || addr.IdUser != userID {
		return nil, errors.New("alamat pengiriman not found or unauthorized")
	}

	trxID := s.idGen.GenerateID()
	kodeInvoice := fmt.Sprintf("INV-%s-%s", time.Now().Format("20060102"), trxID[:6])

	var details []DetailTrx
	var productsToUpdate []produk.Produk
	var grandTotal int

	// 2. Loop through cart items and build logic
	for _, itemReq := range req.Items {
		p, err := s.produkRepo.FindByID(ctx, itemReq.IdProduk)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("product %s not found", itemReq.IdProduk)
		}

		// Verify Stock
		if p.Stok < itemReq.Kuantitas {
			return nil, fmt.Errorf("insufficient stock for product %s (available: %d)", p.NamaProduk, p.Stok)
		}

		// Calculate Subtotal (Assuming normal HargaKonsumen for now)
		subTotal := p.HargaKonsumen * itemReq.Kuantitas
		grandTotal += subTotal

		// Build Detail
		details = append(details, DetailTrx{
			ID:          s.idGen.GenerateID(),
			IdTrx:       trxID,
			IdLogProduk: p.ID, // In a real strict implementation, we'd fetch the latest log_produk ID here. Linking to product ID for now to satisfy foreign constraints based on schema setup if log_produk shares the same ID logic or we just grab the latest log. (For this mini-project, assuming we just log it).
			IdToko:      p.IdToko,
			Kuantitas:   itemReq.Kuantitas,
			HargaTotal:  subTotal,
		})

		// Deduct Stock in memory
		p.Stok -= itemReq.Kuantitas
		productsToUpdate = append(productsToUpdate, *p)
	}

	// 3. Build Transaction Header
	trx := &Trx{
		ID:               trxID,
		IdUser:           userID,
		AlamatPengiriman: req.AlamatPengiriman,
		HargaTotal:       grandTotal,
		KodeInvoice:      kodeInvoice,
		MethodBayar:      req.MethodBayar,
	}

	// 4. Execute massive ACID Transaction
	if err := s.trxRepo.CreateCheckout(ctx, trx, details, productsToUpdate); err != nil {
		return nil, errors.New("checkout failed due to system error, transaction rolled back")
	}

	return mapToTransactionResponse(trx), nil
}

func (s *transactionService) GetMyTransactions(ctx context.Context, userID string) ([]TransactionResponse, error) {
	trxs, err := s.trxRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []TransactionResponse
	for _, t := range trxs {
		res = append(res, *mapToTransactionResponse(&t))
	}
	return res, nil
}

func (s *transactionService) GetTransactionDetail(ctx context.Context, userID, trxID string) (*TransactionResponse, error) {
	// In a complete app, this would also fetch `DetailTrx` and return a nested object.
	// For now, we secure it and return the header.
	trxs, err := s.trxRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var target *Trx
	for _, t := range trxs {
		if t.ID == trxID {
			target = &t
			break
		}
	}

	if target == nil {
		return nil, errors.New("transaction not found or unauthorized")
	}

	return mapToTransactionResponse(target), nil
}

func mapToTransactionResponse(t *Trx) *TransactionResponse {
	return &TransactionResponse{
		ID:               t.ID,
		KodeInvoice:      t.KodeInvoice,
		HargaTotal:       t.HargaTotal,
		MethodBayar:      t.MethodBayar,
		AlamatPengiriman: t.AlamatPengiriman,
		CreatedAt:        t.CreatedAt,
		// Assuming status defaults to "PAID" or "PENDING", can be added to entity later.
	}
}
