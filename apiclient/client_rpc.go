package apiclient

import (
	"context"
	"tgtest/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BalanceClientRPC interface {
	GetBalance(ctx context.Context, accountID string) (*pb.GetBalanceResponse, error)
	Close() error
}

type clientRPC struct {
	client pb.BalanceServiceClient
	conn   *grpc.ClientConn
}

func CreateClientRPC() (BalanceClientRPC, error) {
	conn, err := grpc.NewClient("balance-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &clientRPC{client: pb.NewBalanceServiceClient(conn), conn: conn}, nil
}
func (c *clientRPC) GetBalance(ctx context.Context, accountID string) (*pb.GetBalanceResponse, error) {
	return c.client.GetBalance(ctx, &pb.GetBalanceRequest{AccountId: accountID})
}
func (c *clientRPC) Close() error {
	return c.conn.Close()
}
