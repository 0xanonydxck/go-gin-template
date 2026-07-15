package user

import (
	"context"
	"fmt"

	"github.com/chai-rs/simple-bookstore/infrastructure/auth"
	"github.com/chai-rs/simple-bookstore/internal/model"
	"github.com/chai-rs/simple-bookstore/pkg/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Service represents the user service interface.
type Service interface {
	Login(ctx context.Context, email string, password string) (string, string, error)
	Register(ctx context.Context, user *model.User) (string, string, error)
	Logout(ctx context.Context, metadata *auth.AccessProperties) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

// service implements the Service interface
type service struct {
	repo         Repository
	auth         auth.Auth
	tokenManager auth.TokenManager
	enforcer     auth.AuthEnforcer
}

func NewService(repo Repository, auth auth.Auth, tokenManager auth.TokenManager, enforcer auth.AuthEnforcer) *service {
	return &service{repo, auth, tokenManager, enforcer}
}

func (s *service) Login(ctx context.Context, email string, password string) (string, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("email", email).Msg("🚨 failed to get user by email")
		return "", "", err
	}

	if !crypto.ComparePassword(password, user.HashedPassword) {
		log.Ctx(ctx).Error().Str("email", email).Msg("🚨 invalid password")
		return "", "", fmt.Errorf("invalid password")
	}

	ts, err := s.tokenManager.CreateToken(user.ID.String(), user.Email)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to create token")
		return "", "", err
	}

	if err := s.auth.CreateAuth(ctx, user.ID.String(), ts); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to create auth")
		return "", "", err
	}

	return ts.AccessToken, ts.RefreshToken, nil
}

func (s *service) Register(ctx context.Context, user *model.User) (string, string, error) {
	user.ID = uuid.New()
	err := s.repo.Create(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("email", user.Email).Msg("🚨 failed to create user")
		return "", "", err
	}

	if err := s.enforcer.AddPolicy(user.ID.String(), auth.Resource, auth.Read); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("email", user.Email).Msg("🚨 failed to add policy")
		return "", "", err
	}

	if err := s.enforcer.AddPolicy(user.ID.String(), auth.Resource, auth.Write); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("email", user.Email).Msg("🚨 failed to add policy")
		return "", "", err
	}

	ts, err := s.tokenManager.CreateToken(user.ID.String(), user.Email)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to create token")
		return "", "", err
	}

	if err := s.auth.CreateAuth(ctx, user.ID.String(), ts); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to create auth")
		return "", "", err
	}

	return ts.AccessToken, ts.RefreshToken, nil
}

func (s *service) Logout(ctx context.Context, metadata *auth.AccessProperties) error {
	if err := s.auth.DeleteAccessToken(ctx, metadata); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to delete access token")
		return err
	}

	if err := s.auth.DeleteRefreshToken(ctx, auth.ToRefreshUUID(metadata.TokenUUID, metadata.UserID)); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to delete refresh token")
		return err
	}

	return nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := s.auth.FetchAuth(ctx, refreshToken)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to fetch auth")
		return "", "", err
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to get user by id")
		return "", "", err
	}

	ts, err := s.tokenManager.CreateToken(user.ID.String(), user.Email)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("🚨 failed to create token")
		return "", "", err
	}

	return ts.AccessToken, ts.RefreshToken, nil
}
