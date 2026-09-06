package user

import (
	"context"
	"errors"
	"time"

	"github.com/ahmadammarm/sommerce-mini-project/internal/toko"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/emsifa"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/hash"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/uid"
	"github.com/ahmadammarm/sommerce-mini-project/utils"
	"github.com/go-playground/validator/v10"
)

// UserService defines the business logic contract for the User domain
type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (string, error)
	GetMyProfile(ctx context.Context, userID string) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserResponse, error)
}

type userService struct {
	userRepo   UserRepository
	tokoRepo   toko.TokoRepository
	wilayahAPI emsifa.WilayahProvider
	idGen      uid.IDGenerator
	validator  *validator.Validate
}

// NewUserService is the constructor for Dependency Injection
func NewUserService(ur UserRepository, tr toko.TokoRepository, wp emsifa.WilayahProvider, idg uid.IDGenerator, v *validator.Validate) UserService {
	return &userService{
		userRepo:   ur,
		tokoRepo:   tr,
		wilayahAPI: wp,
		idGen:      idg,
		validator:  v,
	}
}

// Register orchestrates validation, security, and database persistence to create a new user
func (s *userService) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	// 1. Validate incoming JSON payload structure
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// 2. Strict external validation against Emsifa API
	validProv, err := s.wilayahAPI.IsValidProvinsi(ctx, req.IdProvinsi)
	if err != nil {
		return nil, errors.New("failed to connect to regional verification server")
	}
	if !validProv {
		return nil, errors.New("invalid province ID")
	}

	validKota, err := s.wilayahAPI.IsValidKota(ctx, req.IdProvinsi, req.IdKota)
	if err != nil {
		return nil, errors.New("failed to connect to regional verification server")
	}
	if !validKota {
		return nil, errors.New("invalid city ID")
	}

	// 3. Ensure email uniqueness
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// 4. Hash password and generate primary keys
	hashedPassword, err := hash.HashPassword(req.KataSandi)
	if err != nil {
		return nil, errors.New("failed to process password")
	}
	userID := s.idGen.GenerateID()
	parsedDate, _ := time.Parse("2006-01-02", req.TanggalLahir) // Assuming YYYY-MM-DD

	// 5. Construct and Save User
	user := &User{
		ID:           userID,
		Nama:         req.Nama,
		Email:        req.Email,
		KataSandi:    hashedPassword,
		NoTelp:       req.NoTelp,
		TanggalLahir: parsedDate,
		JenisKelamin: req.JenisKelamin,
		Tentang:      req.Tentang,
		Pekerjaan:    req.Pekerjaan,
		IdProvinsi:   req.IdProvinsi,
		IdKota:       req.IdKota,
		IsAdmin:      false, // Force default to false for security
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 6. Execute Rule 6: Automatically create an empty Toko
	tokoBaru := &toko.Toko{
		ID:     s.idGen.GenerateID(),
		IdUser: userID,
	}
	if err := s.tokoRepo.Create(ctx, tokoBaru); err != nil {
		return nil, err
	}

	// 7. Return safe response mapping (excluding password)
	return mapToUserResponse(user), nil
}

// Login validates credentials and issues a JWT token
func (s *userService) Login(ctx context.Context, req *LoginRequest) (string, error) {
	if err := s.validator.Struct(req); err != nil {
		return "", err
	}

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	genericAuthErr := errors.New("invalid email or password")
	if user == nil {
		return "", genericAuthErr
	}

	if !hash.CheckPasswordHash(req.KataSandi, user.KataSandi) {
		return "", genericAuthErr
	}

	return utils.GenerateToken(user.ID, user.IsAdmin)
}

// GetMyProfile fetches the currently logged-in user's data
func (s *userService) GetMyProfile(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return mapToUserResponse(user), nil
}

// UpdateProfile allows a user to update their personal information
func (s *userService) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*UserResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}

	// If region is being updated, validate with Emsifa
	if req.IdProvinsi != "" && req.IdKota != "" {
		validKota, err := s.wilayahAPI.IsValidKota(ctx, req.IdProvinsi, req.IdKota)
		if err != nil {
			return nil, errors.New("failed to connect to regional verification server")
		}
		if !validKota {
			return nil, errors.New("invalid province or city ID")
		}
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Nama != "" {
		user.Nama = req.Nama
	}
	if req.NoTelp != "" {
		user.NoTelp = req.NoTelp
	}
	if req.TanggalLahir != "" {
		parsedDate, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err == nil {
			user.TanggalLahir = parsedDate
		}
	}
	if req.JenisKelamin != "" {
		user.JenisKelamin = req.JenisKelamin
	}
	if req.Tentang != "" {
		user.Tentang = req.Tentang
	}
	if req.Pekerjaan != "" {
		user.Pekerjaan = req.Pekerjaan
	}
	if req.IdProvinsi != "" {
		user.IdProvinsi = req.IdProvinsi
	}
	if req.IdKota != "" {
		user.IdKota = req.IdKota
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return mapToUserResponse(user), nil
}

// mapToUserResponse is a private helper to map the DB entity to the API DTO
func mapToUserResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		Nama:         u.Nama,
		Email:        u.Email,
		NoTelp:       u.NoTelp,
		TanggalLahir: u.TanggalLahir,
		JenisKelamin: u.JenisKelamin,
		Tentang:      u.Tentang,
		Pekerjaan:    u.Pekerjaan,
		IdProvinsi:   u.IdProvinsi,
		IdKota:       u.IdKota,
		IsAdmin:      u.IsAdmin,
		CreatedAt:    u.CreatedAt,
	}
}
