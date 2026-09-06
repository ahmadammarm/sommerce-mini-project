package user

import (
	"context"
	"errors"
	"time"

	"github.com/ahmadammarm/sommerce-mini-project/internal/toko"
	"github.com/ahmadammarm/sommerce-mini-project/pkg/emsifa"
	"github.com/ahmadammarm/sommerce-mini-project/utils"
	"github.com/go-playground/validator/v10"
)

// UserService defines the business logic contract for the User domain
type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (string, error)
}

type userService struct {
	userRepo   UserRepository
	tokoRepo   toko.TokoRepository
	wilayahAPI emsifa.WilayahProvider
	idGen      utils.IDGenerator
	validator  *validator.Validate
}

// NewUserService is the constructor for Dependency Injection
func NewUserService(ur UserRepository, tr toko.TokoRepository, wp emsifa.WilayahProvider, idg utils.IDGenerator, v *validator.Validate) UserService {
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
	hashedPassword, err := utils.HashPassword(req.KataSandi)
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
		// In a fully robust system, you might want to rollback the user here or use a DB transaction
		// However, doing external API calls inside a DB transaction is an anti-pattern.
		// Doing it here sequentially is the acceptable Clean Architecture tradeoff.
		return nil, err
	}

	// 7. Return safe response mapping (excluding password)
	return &UserResponse{
		ID:           user.ID,
		Nama:         user.Nama,
		Email:        user.Email,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir,
		JenisKelamin: user.JenisKelamin,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		IdProvinsi:   user.IdProvinsi,
		IdKota:       user.IdKota,
		IsAdmin:      user.IsAdmin,
		CreatedAt:    user.CreatedAt,
	}, nil
}

// Login validates credentials and issues a JWT token
func (s *userService) Login(ctx context.Context, req *LoginRequest) (string, error) {
	// 1. Validate request structure
	if err := s.validator.Struct(req); err != nil {
		return "", err
	}

	// 2. Fetch User by Email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}
	// Use a generic error message for both cases to prevent user enumeration attacks
	genericAuthErr := errors.New("invalid email or password")
	if user == nil {
		return "", genericAuthErr
	}

	// 3. Verify bcrypt hash
	if !utils.CheckPasswordHash(req.KataSandi, user.KataSandi) {
		return "", genericAuthErr
	}

	// 4. Generate and return JWT Token
	return utils.GenerateToken(user.ID, user.IsAdmin)
}
