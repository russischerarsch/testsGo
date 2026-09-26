package serv

import (
	"context"
	"strings"
	"testing"
	"tgtest/domain"
	"tgtest/serv/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repository := mocks.NewRepoInterface(t)
	repository.On("CreateUser", mock.Anything, mock.Anything).Return("", domain.ErrUserAlreadyExists).Once()
	service, err := CreateServ(repository)
	require.NoError(t, err)
	_, err = service.CreateUser(context.Background(), "Иван", "aarara@gmail.com", "Qwerty123!", "19")
	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)

}

func TestCreateUser_TimeOutExceeded(t *testing.T) {
	repository := mocks.NewRepoInterface(t)
	repository.On("CreateUser", mock.Anything, mock.Anything).Return("0", context.DeadlineExceeded)
	service, err := CreateServ(repository)
	require.NoError(t, err)
	_, err = service.CreateUser(context.Background(), "Oleg", "oleg@example.com", "Qwerty123!", "20")
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
func TestCreateUser_ValidateName(t *testing.T) {
	repository := mocks.NewRepoInterface(t)
	service, err := CreateServ(repository)
	require.NoError(t, err)
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
	service, err := CreateServ(repository)
	require.NoError(t, err)
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
	service, err := CreateServ(repository)
	require.NoError(t, err)
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
