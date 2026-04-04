package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/jvzito/airball/internal/models"
	"github.com/jvzito/airball/internal/repository"
)

type AuthService struct {
	users     *repository.UserRepo
	jwtSecret []byte
}

func NewAuthService(users *repository.UserRepo, secret string) *AuthService {
	return &AuthService{users: users, jwtSecret: []byte(secret)}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := s.users.Create(ctx, email, string(hash))
	if err != nil {
		return nil, errors.New("email ja cadastrado")
	}
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("credenciais invalidas")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("credenciais invalidas")
	}
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}
