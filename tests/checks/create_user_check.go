package checks

import (
	"testing"
	"tgtest/tests/clients"

	"github.com/stretchr/testify/require"
)

func CheckCreateUserResponse(t *testing.T, resp *clients.CreateUserResponse, expectedID string) {
	t.Helper()
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.ID)
	require.Equal(t, expectedID, resp.ID)
}
