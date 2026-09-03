package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// UserRepository defines the strict contract for database operations related to User
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
}

// userRepository is the concrete implementation of the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository is the constructor provided for Google Wire Dependency Injection
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user record into the database
func (r *userRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail searches for a user by their unique email address.
// Returns nil and no error if the user is not found, allowing the service to handle the logic.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User does not exist
		}
		return nil, err // Unexpected DB error
	}
	return &u, nil
}

// FindByID retrieves a user by their CUID.
func (r *userRepository) FindByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// Update saves changes made to an existing user record
func (r *userRepository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
