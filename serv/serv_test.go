package serv

import (
	"context"
	"strings"
	"testing"
	"tgtest/domain"
	"tgtest/serv/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// type repoMock struct {
// 	mock.Mock
// }

//	func (r *repoMock) CreateUser(ctx context.Context, user *domain.User) (string, error) {
//		args := r.Called(ctx, user)
//		return args.String(0), args.Error(1)
//	}
func TestCreateUser_DuplicateEmail(t *testing.T) {
	repository := mocks.NewRepoInterface(t)
	repository.On("CreateUser", mock.Anything, mock.Anything).Return("", domain.ErrUserAlreadyExists).Once()
	service := CreateServ(repository)
	_, err := service.CreateUser(context.Background(), "Иван", "aarara@gmail.com", "Qwerty123!", "19")
	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)

	repository.AssertExpectations(t)
}

func TestCreateUser_ValidateName(t *testing.T) {
	repository := mocks.NewRepoInterface(t)
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
			_, err := service.CreateUser(context.Background(), testCase.input, "vntoebn@gmail.com", "Qwerty123!", "19")
			assert.ErrorIs(t, err, domain.ErrInvalidName)
		})
	}
}
func TestCreateUser_ValidateEmail(t *testing.T) {
	repository := mocks.NewRepoInterface(t)
	repository.On("CreateUser", mock.Anything, mock.Anything)
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
			_, err := service.CreateUser(context.Background(), "Oleg", testCase.input, "Qwerty123!", "19")
			assert.ErrorIs(t, err, domain.ErrInvalidEmail)
		})
	}
}

func TestCreateUser_ValidatePass(t *testing.T) {
	repository := mocks.NewRepoInterface(t)
	service := CreateServ(repository)
	testCases := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "short password",
			input:   "Qwerty!",
			wantErr: domain.ErrShortPassword,
		},
		{
			name:    "Lower letter",
			input:   "qwerty123!",
			wantErr: domain.ErrNoUpperLetter,
		},
		{
			name:    "No special symbols",
			input:   "Qwerty123",
			wantErr: domain.ErrNoSpecialChar,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := service.CreateUser(context.Background(), "Oleg", "oleg@example.com", testCase.input, "19")
			assert.ErrorIs(t, err, testCase.wantErr)
		})
	}
}
