package book

import (
	"context"

	"github.com/chai-rs/simple-bookstore/internal/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Create(ctx context.Context, book *model.Book) error
	Update(ctx context.Context, book *model.Book) error
	GetAll(ctx context.Context) ([]model.Book, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Book, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, book *model.Book) error {
	err := s.repo.Create(ctx, book)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to create book")
		return err
	}

	return nil
}

func (s *service) Update(ctx context.Context, book *model.Book) error {
	err := s.repo.Update(ctx, book)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to update book")
		return err
	}

	return nil
}

func (s *service) GetAll(ctx context.Context) ([]model.Book, error) {
	books, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to get all books")
		return nil, err
	}

	return books, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*model.Book, error) {
	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to get book by id")
		return nil, err
	}

	return book, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to delete book")
		return err
	}

	return nil
}
