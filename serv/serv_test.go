package serv

import (
	"context"
	"errors"
	"strings"
	"testing"
	"tgtest/domain"

	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

func (r *repoMock) CreateUser(ctx context.Context, user *domain.User) (int, error) {
	args := r.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repository := new(repoMock)
	repository.On("CreateUser", mock.Anything, mock.Anything).Return(0, domain.ErrUserAlreadyExists)
	service := CreateServ(repository)
	_, err := service.CreateUser(context.Background(), "Иван", "aarara@gmail.com", "Qwerty123!", 19)
	if err == nil {
		if !errors.Is(err, domain.ErrUserAlreadyExists) {
			t.Fatalf("ожидалась ошибка %v, фактическая ошибка %v", domain.ErrUserAlreadyExists, err)
		}
		t.Fatalf("expected err %v, actual %v", domain.ErrUserAlreadyExists, err)
	}
	repository.AssertExpectations(t)
}

func TestCreateUser_ValidateName(t *testing.T) {
	repository := new(repoMock)
	repository.On("CreateUser", mock.Anything, mock.Anything)
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
			_, err := service.CreateUser(context.Background(), testCase.input, "vntoebn@gmail.com", "Qwerty123!", 19)
			if !errors.Is(err, domain.ErrInvalidName) {
				t.Fatalf("expected error %v, actual error %v", domain.ErrInvalidName, err)
			}
		})
	}
}
func TestCreateUser_ValidateEmail(t *testing.T) {
	repository := new(repoMock)
	repository.On("CreateUser", mock.Anything, mock.Anything)
	service := CreateServ(repository)

	testCases := []struct {
		name   string
		input  string
		caseID int64
	}{
		{
			name:  "empty email",
			input: "",
		},
		{
			name:   "long email",
			input:  strings.Repeat("a", 65) + "@x.com",
			caseID: 11,
		},
		{
			name:   "email without @",
			input:  "jroinboie.com",
			caseID: 12,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := service.CreateUser(context.Background(), "Oleg", testCase.input, "Qwerty123!", 19)
			if !errors.Is(err, domain.ErrInvalidEmail) {
				t.Fatalf("expected error %v, actual %v", domain.ErrInvalidEmail, err)
			}
		})
	}
}
