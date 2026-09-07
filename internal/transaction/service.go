package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ahmadammarm/sommerce-mini-project/internal/alamat"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/uid"
	"github.com/go-playground/validator/v10"
)

type TransactionService interface {
	Checkout(ctx context.Context, userID, idempotencyKey string, req *CheckoutRequest) (*TransactionResponse, error)
	GetMyTransactions(ctx context.Context, userID string) ([]TransactionResponse, error)
	GetTransactionDetail(ctx context.Context, userID, trxID string) (*TransactionDetailResponse, error)
}

type transactionService struct {
	trxRepo    TransactionRepository
	alamatRepo alamat.AlamatRepository
	idGen      uid.IDGenerator
	validator  *validator.Validate
}

func NewTransactionService(tr TransactionRepository, ar alamat.AlamatRepository, idg uid.IDGenerator, v *validator.Validate) TransactionService {
	return &transactionService{
		trxRepo:    tr,
		alamatRepo: ar,
		idGen:      idg,
		validator:  v,
	}
}

func (s *transactionService) Checkout(ctx context.Context, userID, idempotencyKey string, req *CheckoutRequest) (*TransactionResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	addr, err := s.alamatRepo.FindByID(ctx, req.AlamatPengiriman)
	if err != nil {
		return nil, err
	}
	if addr == nil || addr.IdUser != userID {
		return nil, errors.New("alamat pengiriman not found or unauthorized")
	}

	trxID := s.idGen.GenerateID()
	kodeInvoice := fmt.Sprintf("INV-%s-%s", time.Now().Format("20060102"), trxID[:6])

	trx := &Trx{
		ID:               trxID,
		IdUser:           userID,
		AlamatPengiriman: req.AlamatPengiriman,
		IdempotencyKey:   idempotencyKey,
		KodeInvoice:      kodeInvoice,
		MethodBayar:      req.MethodBayar,
	}

	if err := s.trxRepo.CreateCheckout(ctx, trx, req.Items); err != nil {
		return nil, fmt.Errorf("checkout failed: %v", err)
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

func (s *transactionService) GetTransactionDetail(ctx context.Context, userID, trxID string) (*TransactionDetailResponse, error) {
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

	details, err := s.trxRepo.FindDetailsByTrxID(ctx, trxID)
	if err != nil {
		return nil, err
	}

	var detailsResp []DetailTransactionResponse
	for _, d := range details {
		detailsResp = append(detailsResp, DetailTransactionResponse{
			ID:         d.ID,
			Kuantitas:  d.Kuantitas,
			HargaTotal: d.HargaTotal,
			Produk: LogProdukResponse{
				ID:            d.LogProduk.ID,
				IdProduk:      d.LogProduk.IdProduk,
				NamaProduk:    d.LogProduk.NamaProduk,
				HargaReseller: d.LogProduk.HargaReseller,
				HargaKonsumen: d.LogProduk.HargaKonsumen,
				Deskripsi:     d.LogProduk.Deskripsi,
			},
		})
	}

	resp := &TransactionDetailResponse{
		TransactionResponse: *mapToTransactionResponse(target),
		Details:             detailsResp,
	}

	return resp, nil
}

func mapToTransactionResponse(t *Trx) *TransactionResponse {
	return &TransactionResponse{
		ID:               t.ID,
		KodeInvoice:      t.KodeInvoice,
		HargaTotal:       t.HargaTotal,
		MethodBayar:      t.MethodBayar,
		AlamatPengiriman: t.AlamatPengiriman,
		CreatedAt:        t.CreatedAt,
	}
}
