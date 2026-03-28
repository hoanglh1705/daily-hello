package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo        *repositories.UserRepository
	tokenRepo       repositories.TokenRepository
	jwtSecret       string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	tokenRepo repositories.TokenRepository,
	jwtSecret string,
	accessDurationSec int,
	refreshDurationSec int,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		jwtSecret:       jwtSecret,
		accessDuration:  time.Duration(accessDurationSec) * time.Second,
		refreshDuration: time.Duration(refreshDurationSec) * time.Second,
	}
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, appErrors.ErrInvalidCreds
	}

	if user.Status != "active" {
		return nil, appErrors.ErrAccountInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, appErrors.ErrInvalidCreds
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	refreshToken, err := s.createRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessDuration.Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req models.LogoutRequest) error {
	return s.tokenRepo.DeleteByToken(ctx, req.RefreshToken)
}

func (s *AuthService) RefreshToken(ctx context.Context, req models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	rt, err := s.tokenRepo.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, appErrors.ErrInvalidToken
	}

	if rt.IsExpired() {
		_ = s.tokenRepo.DeleteByToken(ctx, req.RefreshToken)
		return nil, appErrors.ErrTokenExpired
	}

	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, appErrors.ErrInvalidToken
	}

	if user.Status != "active" {
		return nil, appErrors.ErrAccountInactive
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	return &models.RefreshTokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.accessDuration.Seconds()),
	}, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    string(user.Role),
		"exp":     time.Now().Add(s.accessDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) createRefreshToken(ctx context.Context, userID uint) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(raw)

	rt := &models.RefreshToken{
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(s.refreshDuration),
	}

	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return "", err
	}

	return tokenStr, nil
}
