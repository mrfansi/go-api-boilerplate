package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/errors"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/repository"
	"github.com/mrfansi/go-api-boilerplate/internal/infrastructure/config"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	RefreshToken(token string) (string, error)
}

type authService struct {
	config   *config.Config
	userRepo repository.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo repository.UserRepository) AuthService {
	return &authService{
		config:   cfg,
		userRepo: userRepo,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.ErrInvalidCredential
	}

	if err := user.ComparePassword(password); err != nil {
		return "", errors.ErrInvalidCredential
	}

	if !user.Active {
		return "", errors.ErrUnauthorized
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.config.JWT.ExpirationHours).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, errors.ErrTokenExpired
		}
		return nil, errors.ErrInvalidToken
	}

	if !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	return token, nil
}

func (s *authService) RefreshToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil && err != errors.ErrTokenExpired {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.ErrInvalidToken
	}

	// Create new token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims["sub"],
		"email": claims["email"],
		"role":  claims["role"],
		"exp":   time.Now().Add(s.config.JWT.ExpirationHours).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err = newToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
