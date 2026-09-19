package serv

import (
	"context"
	"errors"
	"testing"
	"tgtest/domain"
)

type repoMock struct {
	createFunc func(ctx context.Context, user *domain.User) (int, error)
}

func (r repoMock) CreateUser(ctx context.Context, user *domain.User) (int, error) {
	return r.createFunc(ctx, user)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repository := repoMock{
		createFunc: func(ctx context.Context, user *domain.User) (int, error) {
			return 0, domain.ErrUserAlreadyExists
		},
	}
	service := CreateServ(repository)
	_, err := service.CreateUser(context.Background(), "Иван", "aarara@gmail.com")
	if err != nil {
		if !errors.Is(err, domain.ErrUserAlreadyExists) {
			t.Fatal("ожидалась ошибка %w, фактическая ошибка %w", domain.ErrUserAlreadyExists, err)
		}
	}
}
