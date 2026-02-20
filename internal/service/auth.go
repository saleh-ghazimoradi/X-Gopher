package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/config"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
	"github.com/saleh-ghazimoradi/X-Gopher/utils"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, input *dto.RegisterReq) (*dto.AuthResp, error)
	Login(ctx context.Context, input *dto.LoginReq) (*dto.AuthResp, error)
	RefreshToken(ctx context.Context, input *dto.RefreshTokenReq) (*dto.AuthResp, error)
	Logout(ctx context.Context, input *dto.RefreshTokenReq) error
}

type authService struct {
	config          *config.Config
	userRepository  repository.UserRepository
	tokenRepository repository.TokenRepository
}

func (a *authService) Register(ctx context.Context, input *dto.RegisterReq) (*dto.AuthResp, error) {
	existing, err := a.userRepository.GetUserByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if existing != nil {
		return nil, repository.ErrDuplicateEmail
	}

	user, err := a.toUser(input)
	if err != nil {
		return nil, err
	}

	if err := a.userRepository.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return a.generateAuthResp(ctx, user)
}

func (a *authService) Login(ctx context.Context, input *dto.LoginReq) (*dto.AuthResp, error) {
	user, err := a.userRepository.GetUserByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
	}

	if !utils.CheckPasswordHash(user.Password, input.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return a.generateAuthResp(ctx, user)
}

func (a *authService) RefreshToken(ctx context.Context, input *dto.RefreshTokenReq) (*dto.AuthResp, error) {
	claim, err := utils.ValidateToken(input.RefreshToken, a.config.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	refreshToken, err := a.tokenRepository.GetValidRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	user, err := a.userRepository.GetUserById(ctx, claim.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	if err := a.tokenRepository.DeleteRefreshTokenById(ctx, refreshToken.Id); err != nil {
		return nil, fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return a.generateAuthResp(ctx, user)
}

func (a *authService) Logout(ctx context.Context, input *dto.RefreshTokenReq) error {
	return a.tokenRepository.DeleteRefreshToken(ctx, input.RefreshToken)
}

func (a *authService) toUser(input *dto.RegisterReq) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
		Followers: make([]string, 0),
		Following: make([]string, 0),
	}, nil
}

func (a *authService) generateAuthResp(ctx context.Context, user *domain.User) (*dto.AuthResp, error) {
	accessToken, refreshToken, err := utils.GenerateToken(a.config, user.Id, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refresh := &domain.RefreshToken{
		UserId:    user.Id,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(a.config.JWT.RefreshTokenExpires),
		CreatedAt: time.Now(),
	}

	if err := a.tokenRepository.CreateRefreshToken(ctx, refresh); err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &dto.AuthResp{
		User: dto.UserResp{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			ImageUrl:  user.ImageUrl,
			Followers: user.Followers,
			Following: user.Following,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewAuthService(config *config.Config, userRepository repository.UserRepository, tokenRepository repository.TokenRepository) AuthService {
	return &authService{
		config:          config,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}
