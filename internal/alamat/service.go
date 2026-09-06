package alamat

import (
	"context"
	"errors"

	"github.com/ahmadammarm/sommerce-mini-project/pkg/emsifa"
	"github.com/ahmadammarm/sommerce-mini-project/utils"
	"github.com/go-playground/validator/v10"
)

type AlamatService interface {
	CreateAlamat(ctx context.Context, userID string, req *CreateAlamatRequest) (*AlamatResponse, error)
	GetMyAlamat(ctx context.Context, userID string) ([]AlamatResponse, error)
	GetAlamatByID(ctx context.Context, userID, alamatID string) (*AlamatResponse, error)
	UpdateAlamat(ctx context.Context, userID, alamatID string, req *UpdateAlamatRequest) (*AlamatResponse, error)
	DeleteAlamat(ctx context.Context, userID, alamatID string) error
}

type alamatService struct {
	repo       AlamatRepository
	wilayahAPI emsifa.WilayahProvider
	idGen      utils.IDGenerator
	validator  *validator.Validate
}

func NewAlamatService(r AlamatRepository, wp emsifa.WilayahProvider, idg utils.IDGenerator, v *validator.Validate) AlamatService {
	return &alamatService{repo: r, wilayahAPI: wp, idGen: idg, validator: v}
}

func (s *alamatService) CreateAlamat(ctx context.Context, userID string, req *CreateAlamatRequest) (*AlamatResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// External Validation
	validProv, err := s.wilayahAPI.IsValidProvinsi(ctx, req.IdProvinsi)
	if err != nil || !validProv {
		return nil, errors.New("invalid province ID")
	}

	validKota, err := s.wilayahAPI.IsValidKota(ctx, req.IdProvinsi, req.IdKota)
	if err != nil || !validKota {
		return nil, errors.New("invalid city ID")
	}

	alamat := &Alamat{
		ID:           s.idGen.GenerateID(),
		IdUser:       userID,
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		DetailAlamat: req.DetailAlamat,
		IdProvinsi:   req.IdProvinsi,
		IdKota:       req.IdKota,
	}

	if err := s.repo.Create(ctx, alamat); err != nil {
		return nil, err
	}
	return mapToAlamatResponse(alamat), nil
}

func (s *alamatService) GetMyAlamat(ctx context.Context, userID string) ([]AlamatResponse, error) {
	alamats, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []AlamatResponse
	for _, a := range alamats {
		responses = append(responses, *mapToAlamatResponse(&a))
	}
	return responses, nil
}

func (s *alamatService) GetAlamatByID(ctx context.Context, userID, alamatID string) (*AlamatResponse, error) {
	alamat, err := s.repo.FindByID(ctx, alamatID)
	if err != nil {
		return nil, err
	}

	// Data Isolation Security Check
	if alamat == nil || alamat.IdUser != userID {
		return nil, errors.New("address not found or unauthorized")
	}

	return mapToAlamatResponse(alamat), nil
}

func (s *alamatService) UpdateAlamat(ctx context.Context, userID, alamatID string, req *UpdateAlamatRequest) (*AlamatResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Verify Ownership
	alamat, err := s.repo.FindByID(ctx, alamatID)
	if err != nil {
		return nil, err
	}
	if alamat == nil || alamat.IdUser != userID {
		return nil, errors.New("address not found or unauthorized")
	}

	// Emsifa Validation if region is updated
	if req.IdProvinsi != "" && req.IdKota != "" {
		validKota, err := s.wilayahAPI.IsValidKota(ctx, req.IdProvinsi, req.IdKota)
		if err != nil || !validKota {
			return nil, errors.New("invalid province or city ID")
		}
	} else if req.IdProvinsi != "" || req.IdKota != "" {
		return nil, errors.New("both province and city ID must be provided together for update")
	}

	// Update fields
	if req.JudulAlamat != "" {
		alamat.JudulAlamat = req.JudulAlamat
	}
	if req.NamaPenerima != "" {
		alamat.NamaPenerima = req.NamaPenerima
	}
	if req.NoTelp != "" {
		alamat.NoTelp = req.NoTelp
	}
	if req.DetailAlamat != "" {
		alamat.DetailAlamat = req.DetailAlamat
	}
	if req.IdProvinsi != "" {
		alamat.IdProvinsi = req.IdProvinsi
	}
	if req.IdKota != "" {
		alamat.IdKota = req.IdKota
	}

	if err := s.repo.Update(ctx, alamat); err != nil {
		return nil, err
	}
	return mapToAlamatResponse(alamat), nil
}

func (s *alamatService) DeleteAlamat(ctx context.Context, userID, alamatID string) error {
	alamat, err := s.repo.FindByID(ctx, alamatID)
	if err != nil {
		return err
	}
	if alamat == nil || alamat.IdUser != userID {
		return errors.New("address not found or unauthorized")
	}

	return s.repo.Delete(ctx, alamatID)
}

func mapToAlamatResponse(a *Alamat) *AlamatResponse {
	return &AlamatResponse{
		ID:           a.ID,
		JudulAlamat:  a.JudulAlamat,
		NamaPenerima: a.NamaPenerima,
		NoTelp:       a.NoTelp,
		DetailAlamat: a.DetailAlamat,
		IdProvinsi:   a.IdProvinsi,
		IdKota:       a.IdKota,
	}
}
