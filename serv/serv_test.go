package serv

import (
	"context"
	"errors"
	"strings"
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
			t.Fatalf("ожидалась ошибка %v, фактическая ошибка %v", domain.ErrUserAlreadyExists, err)
		}
	}
}

func TestCreateUser_ValidateName(t *testing.T) {
	repository := repoMock{
		createFunc: func(ctx context.Context, user *domain.User) (int, error) {
			t.Fatal("repo is not supposed to be called while validating input")
			return 0, nil
		},
	}
	service := CreateServ(repository)
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "empty name",
			input: "",
		},

		{
			name:  "only spaces",
			input: " 		"},

		{
			name:  "one symbol",
			input: "g",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := service.CreateUser(context.Background(), testCase.input, "vntoebn@gmail.com")
			if !errors.Is(err, domain.ErrInvalidName) {
				t.Fatalf("expected error %v, actual error %v", domain.ErrInvalidName, err)
			}
		})
	}
}
func TestCreateUser_ValidateEmail(t *testing.T) {
	repository := repoMock{createFunc: func(ctx context.Context, user *domain.User) (int, error) {
		t.Fatal("repo is not supposed to be called")
		return 0, nil
	}}
	service := CreateServ(repository)

	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "empty email",
			input: "",
		},
		{
			name:  "long email",
			input: strings.Repeat("a", 65) + "@x.com",
		},
		{
			name:  "email without @",
			input: "jroinboie.com",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := service.CreateUser(context.Background(), "Oleg", testCase.input)
			if !errors.Is(err, domain.ErrInvalidEmail) {
				t.Fatalf("expected error %v, actual %v", domain.ErrInvalidEmail, err)
			}
		})
	}
}
