package toko

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
)

type TokoService interface {
	GetMyToko(ctx context.Context, userID string) (*TokoResponse, error)
	UpdateMyToko(ctx context.Context, userID string, req *UpdateTokoRequest) (*TokoResponse, error)
	GetTokoByID(ctx context.Context, tokoID string) (*TokoResponse, error)
}

type tokoService struct {
	repo      TokoRepository
	validator *validator.Validate
}

func NewTokoService(r TokoRepository, v *validator.Validate) TokoService {
	return &tokoService{repo: r, validator: v}
}

// GetMyToko retrieves the store created during user registration
func (s *tokoService) GetMyToko(ctx context.Context, userID string) (*TokoResponse, error) {
	toko, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if toko == nil {
		return nil, errors.New("store not found")
	}

	return mapToTokoResponse(toko), nil
}

// UpdateMyToko lets the user change their store name and photo
func (s *tokoService) UpdateMyToko(ctx context.Context, userID string, req *UpdateTokoRequest) (*TokoResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	toko, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if toko == nil {
		return nil, errors.New("store not found")
	}

	if req.NamaToko != "" {
		toko.NamaToko = req.NamaToko
	}
	if req.UrlFoto != "" {
		toko.UrlFoto = req.UrlFoto
	}

	if err := s.repo.Update(ctx, toko); err != nil {
		return nil, err
	}

	return mapToTokoResponse(toko), nil
}

// GetTokoByID fetches public store information
func (s *tokoService) GetTokoByID(ctx context.Context, tokoID string) (*TokoResponse, error) {
	toko, err := s.repo.FindByID(ctx, tokoID)
	if err != nil {
		return nil, err
	}
	if toko == nil {
		return nil, errors.New("store not found")
	}

	return mapToTokoResponse(toko), nil
}

func mapToTokoResponse(t *Toko) *TokoResponse {
	return &TokoResponse{
		ID:       t.ID,
		NamaToko: t.NamaToko,
		UrlFoto:  t.UrlFoto,
	}
}
