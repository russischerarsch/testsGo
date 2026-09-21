package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestCreateUserE2E(t *testing.T) {
	reqBody := map[string]string{
		"name":  "Oleg",
		"email": "oleg@example.com",
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
}
