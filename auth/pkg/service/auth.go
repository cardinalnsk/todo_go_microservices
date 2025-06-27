package service

import (
	"auth/model"
	"auth/pkg/config"
	"auth/pkg/repository"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	bcryptCost = 12
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
	cfg  config.AuthConfig
}

func NewAuthService(repo repository.Authorization, cfg config.AuthConfig) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	pwdHash, err := generatePasswordHash(user.Password + s.cfg.Salt)
	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}
	user.Password = pwdHash
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}

	if err := comparePasswordHash(user.Password, password+s.cfg.Salt); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Id,
	})

	return token.SignedString(s.cfg.PrivateKey)
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		// Можно вернуть либо приватный, либо публичный ключ:
		return s.cfg.PrivateKey, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func comparePasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
