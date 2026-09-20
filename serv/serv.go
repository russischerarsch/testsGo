package serv

import (
	"context"
	"strconv"
	"strings"
	"tgtest/domain"
)

type RepoInterface interface {
	CreateUser(ctx context.Context, user *domain.User) (int, error)
}

type Service struct {
	repo RepoInterface
}

func CreateServ(repo RepoInterface) *Service {
	return &Service{repo: repo}
}
func (s *Service) CreateUser(ctx context.Context, name, email string) (string, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", domain.ErrInvalidName
	}
	if len(name) < 2 || len(name) > 70 {
		return "", domain.ErrInvalidName
	}
	if !strings.Contains(email, "@") || len(email) > 70 {
		return "", domain.ErrInvalidEmail
	}
	var user = &domain.User{
		Name:  name,
		Email: email,
	}
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(id), nil
}
