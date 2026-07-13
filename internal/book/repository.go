package book

import (
	"context"

	errs "github.com/chai-rs/simple-bookstore/internal/error"
	"github.com/chai-rs/simple-bookstore/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) ([]model.Book, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Book, error)
	Create(ctx context.Context, book *model.Book) error
	Update(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, book *model.Book) error {
	if err := r.db.WithContext(ctx).Create(book).Error; err != nil {
		return errs.FromGorm(err)
	}

	return nil
}

func (r *repository) GetAll(ctx context.Context) ([]model.Book, error) {
	var books []model.Book
	if err := r.db.WithContext(ctx).Preload("Genre").Preload("Tags").Find(&books).Error; err != nil {
		return nil, errs.FromGorm(err)
	}

	return books, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*model.Book, error) {
	var book model.Book
	if err := r.db.WithContext(ctx).Preload("Genre").Preload("Tags").Where("id = ?", id).First(&book).Error; err != nil {
		return nil, errs.FromGorm(err)
	}

	return &book, nil
}

func (r *repository) Update(ctx context.Context, book *model.Book) error {
	if err := r.db.WithContext(ctx).Save(book).Error; err != nil {
		return errs.FromGorm(err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(&model.Book{}, "id = ?", id).Error; err != nil {
		return errs.FromGorm(err)
	}

	return nil
}
