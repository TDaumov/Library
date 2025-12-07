package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"library_vebservice/models"
	"library_vebservice/repository"

	//"github.com/yourname/library-api/internal/models"
	//"github.com/yourname/library-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, email, password, role string) (*models.User, error)
	Login(ctx context.Context, email, password, secret string) (string, error)
}

type authService struct {
	users repository.UserRepo
}

func NewAuthService(u repository.UserRepo) AuthService {
	return &authService{users: u}
}

func (s *authService) Register(ctx context.Context, email, password, role string) (*models.User, error) {
	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		Email:    email,
		Password: string(hashed),
		Role:     role,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password, secret string) (string, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	claims := jwt.MapClaims{
		"sub":   u.ID.String(),
		"email": u.Email,
		"role":  u.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
