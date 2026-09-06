package category

import (
	"context"
	"errors"

	"github.com/ahmadammarm/sommerce-mini-project/pkg/uid"
	"github.com/go-playground/validator/v10"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *CategoryRequest) (*CategoryResponse, error)
	GetAllCategories(ctx context.Context) ([]CategoryResponse, error)
	GetCategoryByID(ctx context.Context, id string) (*CategoryResponse, error)
	UpdateCategory(ctx context.Context, id string, req *CategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id string) error
}

type categoryService struct {
	repo      CategoryRepository
	idGen     uid.IDGenerator
	validator *validator.Validate
}

func NewCategoryService(r CategoryRepository, idg uid.IDGenerator, v *validator.Validate) CategoryService {
	return &categoryService{repo: r, idGen: idg, validator: v}
}

func (s *categoryService) CreateCategory(ctx context.Context, req *CategoryRequest) (*CategoryResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	category := &Category{
		ID:           s.idGen.GenerateID(),
		NamaCategory: req.NamaCategory,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return mapToCategoryResponse(category), nil
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]CategoryResponse, error) {
	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []CategoryResponse
	for _, c := range categories {
		responses = append(responses, *mapToCategoryResponse(&c))
	}
	return responses, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*CategoryResponse, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	return mapToCategoryResponse(category), nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req *CategoryRequest) (*CategoryResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	category.NamaCategory = req.NamaCategory

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return mapToCategoryResponse(category), nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	return s.repo.Delete(ctx, id)
}

func mapToCategoryResponse(c *Category) *CategoryResponse {
	return &CategoryResponse{
		ID:           c.ID,
		NamaCategory: c.NamaCategory,
	}
}
