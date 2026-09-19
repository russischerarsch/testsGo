package serv

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"tgtest/domain"
)

var ErrInvalidName = errors.New("невалидное имя")

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
		return "", ErrInvalidName
	}
	if len(name) < 2 || len(name) > 70 {
		return "", ErrInvalidName
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
