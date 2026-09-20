package serv

import (
	"context"
	"testing"
	"tgtest/repo"
	"tgtest/serv"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCreateUser_Integration_Success(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container %v", err)
	}
	defer pgContainer.Terminate(ctx)
	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	dsn := "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	pool, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to db %v", err)
	}
	defer pool.Close(ctx)

	repo := repo.CreateRepo(pool)
	svc := serv.CreateServ(repo)
	id, err := svc.CreateUser(ctx, "Oleg", "oleg@exmp.com")

}
