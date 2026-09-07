package produk

import (
	"context"
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/ahmadammarm/sommerce-mini-project/internal/toko"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/response"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/uid"
	"github.com/go-playground/validator/v10"
)

type ProdukService interface {
	CreateProduk(ctx context.Context, userID string, req *CreateProdukRequest) (*ProdukResponse, error)
	UpdateProduk(ctx context.Context, userID, produkID string, req *UpdateProdukRequest) (*ProdukResponse, error)
	GetProdukByID(ctx context.Context, produkID string) (*ProdukResponse, error)
	GetAllProduk(ctx context.Context, filter ProductFilterDTO) (*response.PagedResponse, error)
}

type produkService struct {
	produkRepo ProdukRepository
	tokoRepo   toko.TokoRepository
	idGen      uid.IDGenerator
	validator  *validator.Validate
}

func NewProdukService(pr ProdukRepository, tr toko.TokoRepository, idg uid.IDGenerator, v *validator.Validate) ProdukService {
	return &produkService{
		produkRepo: pr,
		tokoRepo:   tr,
		idGen:      idg,
		validator:  v,
	}
}

func (s *produkService) CreateProduk(ctx context.Context, userID string, req *CreateProdukRequest) (*ProdukResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Rule 6: Get Toko implicitly based on UserID
	vendorToko, err := s.tokoRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if vendorToko == nil {
		return nil, errors.New("unauthorized: vendor store not found")
	}

	produkID := s.idGen.GenerateID()
	slug := generateSlug(req.NamaProduk)

	// 1. Create Main Product
	p := &Produk{
		ID:            produkID,
		NamaProduk:    req.NamaProduk,
		Slug:          slug,
		HargaReseller: req.HargaReseller,
		HargaKonsumen: req.HargaKonsumen,
		Stok:          req.Stok,
		Deskripsi:     req.Deskripsi,
		IdCategory:    req.IdCategory,
		IdToko:        vendorToko.ID,
	}

	// 2. Map Photos
	var fotos []FotoProduk
	for _, url := range req.Photos {
		fotos = append(fotos, FotoProduk{
			ID:       s.idGen.GenerateID(),
			IdProduk: produkID,
			Url:      url,
		})
	}

	if err := s.produkRepo.Create(ctx, p, fotos); err != nil {
		return nil, err
	}

	// Just for returning immediately after creation, attach manually
	p.Fotos = fotos
	return mapToProdukResponse(p), nil
}

func (s *produkService) UpdateProduk(ctx context.Context, userID, produkID string, req *UpdateProdukRequest) (*ProdukResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// Check Vendor Ownership
	vendorToko, err := s.tokoRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if vendorToko == nil {
		return nil, errors.New("unauthorized: vendor store not found")
	}

	// Fetch Existing Product
	p, err := s.produkRepo.FindByID(ctx, produkID)
	if err != nil {
		return nil, err
	}
	if p == nil || p.IdToko != vendorToko.ID {
		return nil, errors.New("product not found or unauthorized")
	}

	// Update fields
	if req.NamaProduk != "" {
		p.NamaProduk = req.NamaProduk
		p.Slug = generateSlug(req.NamaProduk)
	}
	if req.HargaReseller > 0 {
		p.HargaReseller = req.HargaReseller
	}
	if req.HargaKonsumen > 0 {
		p.HargaKonsumen = req.HargaKonsumen
	}
	if req.Stok >= 0 {
		p.Stok = req.Stok
	}
	if req.Deskripsi != "" {
		p.Deskripsi = req.Deskripsi
	}
	if req.IdCategory != "" {
		p.IdCategory = req.IdCategory
	}

	if err := s.produkRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	return mapToProdukResponse(p), nil
}

func (s *produkService) GetProdukByID(ctx context.Context, produkID string) (*ProdukResponse, error) {
	p, err := s.produkRepo.FindByID(ctx, produkID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("product not found")
	}

	return mapToProdukResponse(p), nil
}

func (s *produkService) GetAllProduk(ctx context.Context, filter ProductFilterDTO) (*response.PagedResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	products, totalItems, err := s.produkRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	var data []ProdukResponse
	for _, p := range products {
		data = append(data, *mapToProdukResponse(&p))
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(filter.Limit)))

	return &response.PagedResponse{
		Data: data,
		Meta: response.PaginationMeta{
			CurrentPage: filter.Page,
			Limit:       filter.Limit,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

// Helpers
func mapToProdukResponse(p *Produk) *ProdukResponse {
	var urls []string
	for _, f := range p.Fotos {
		urls = append(urls, f.Url)
	}
	return &ProdukResponse{
		ID:            p.ID,
		NamaProduk:    p.NamaProduk,
		Slug:          p.Slug,
		HargaReseller: p.HargaReseller,
		HargaKonsumen: p.HargaKonsumen,
		Stok:          p.Stok,
		Deskripsi:     p.Deskripsi,
		IdCategory:    p.IdCategory,
		IdToko:        p.IdToko,
		Photos:        urls,
	}
}

func generateSlug(title string) string {
	// Simple slug generator: lowercasing and replacing non-alphanumerics with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := re.ReplaceAllString(strings.ToLower(title), "-")
	return strings.Trim(slug, "-")
}
