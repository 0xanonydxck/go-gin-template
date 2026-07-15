package user

import (
	"context"

	errs "github.com/chai-rs/simple-bookstore/internal/error"
	"github.com/chai-rs/simple-bookstore/internal/model"
	"gorm.io/gorm"
)

// Repository represents the user repository interface.
type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}

// repository implements the Repository interface.
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return errs.FromGorm(err)
	}
	return nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errs.FromGorm(err)
	}

	return &user, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, errs.FromGorm(err)
	}

	return &user, nil
}
