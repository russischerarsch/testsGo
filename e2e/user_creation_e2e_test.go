package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"testing"
	"tgtest/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestCreateUserE2E(t *testing.T) {
	reqBody := map[string]string{
		"name":     "Oleg",
		"email":    "oleg@example.com",
		"password": "Qwerty123!",
		"Age":      "22",
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request body, %v", err)
	}
	resp, err := http.Post(serverURL+"/users", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to send request %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, actual code %v", resp.StatusCode)
	}
	var response struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to get response, %v", err)
	}
	id, err := strconv.Atoi(response.ID)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if id == 0 {
		t.Fatalf("expected non-zero id, got %v", response.ID)
	}
	var name string
	query := `SELECT name FROM users WHERE id = $1`
	if err := testConn.QueryRow(ctx, query, response.ID).Scan(&name); err != nil {
		t.Fatalf("failed to send query to database, %v", err)
	}
	if name != "Oleg" {
		t.Fatalf("expected name 'Oleg', got %v", name)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		query := `
		DELETE FROM users
		WHERE id = $1
		`
		if _, err := testConn.Exec(ctx, query, id); err != nil {
			t.Fatalf("failed to delete row")
		}
	})
}
func TestGetBalance_ClientRPC(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterBalanceServiceServer(server, &myServer{})
	go func() {
		if err := server.Serve(lis); err != nil {
			t.Logf("server exited, %v", err)
		}
	}()
	defer server.Stop()
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet, %v", err)
	}
	defer conn.Close()
	client := pb.NewBalanceServiceClient(conn)

	resp, err := client.GetBalance(context.Background(), &pb.GetBalanceRequest{})
	if err != nil {
		t.Fatalf("GetBalance failed: %v", err)
	}

	if resp.Balance != 100 {
		t.Errorf("expected balance 100, got %d", resp.Balance)
	}
}
