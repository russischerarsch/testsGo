//go:build integration

package serv

import (
	"testing"
	"tgtest/repo"
)

func TestCreateUser_Integration_Success(t *testing.T) {
	repo := repo.CreateRepo(testConn)
	svc := CreateServ(repo)
	id, err := svc.CreateUser(ctx, "Oleg", "oleg@exmpl.com", "Qwerty123!", "22")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	var name string
	query := `
	SELECT name FROM users
	WHERE id = $1
	`
	if err := testConn.QueryRow(ctx, query, id).Scan(&name); err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if name != "Oleg" {
		t.Fatalf("expected name 'Oleg', got %v", name)
	}
	t.Cleanup(func() { testConn.Exec(ctx, "DELETE FROM users WHERE id = $1", id) })
}
