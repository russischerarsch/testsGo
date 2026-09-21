package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func CreateUser_TestE2E(t *testing.T) {
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
}
