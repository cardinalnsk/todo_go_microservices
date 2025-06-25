package service

import (
	"auth/model"
	"auth/pkg/config"
	"auth/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repository *repository.Repository, cfg config.AuthConfig) *Service {
	return &Service{
		Authorization: NewAuthService(repository.Authorization, cfg),
	}
}
