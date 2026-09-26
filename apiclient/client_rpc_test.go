//go:build clientgo

package apiclient

import (
	"context"
	"fmt"
	"testing"
	"tgtest/apiclient/mocks"
	"tgtest/proto/pb"

	"github.com/stretchr/testify/assert"
	mockanything "github.com/stretchr/testify/mock"
	"google.golang.org/grpc/status"
)

func TestGetBalance_ClientRPC(t *testing.T) {
	mock := mocks.NewBalanceClientRPC(t)
	mock.On("GetBalance", mockanything.Anything, "1").Return(&pb.GetBalanceResponse{AccountId: "1", Balance: 12000}, nil)
	resp, err := mock.GetBalance(context.Background(), "1")
	fmt.Println(status.Code(err))
	assert.NoError(t, err)
	assert.Equal(t, int64(12000), resp.Balance)
}
