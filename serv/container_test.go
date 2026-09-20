package serv

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testConn *pgx.Conn
var ctx = context.Background()

func SetUp_TestContainers(ctx context.Context) (*pgx.Conn, testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}
	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, pgContainer, err
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, pgContainer, err
	}
	url := "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	pool, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, pgContainer, err
	}
	m, err := migrate.New("file://../migrations", url)
	if err != nil {
		return nil, pgContainer, err
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, pgContainer, err
	}
	return pool, pgContainer, nil
}

func TestMain(m *testing.M) {
	conn, container, err := SetUp_TestContainers(ctx)
	if err != nil {
		log.Fatalf("failed to set up test container %v", err)
	}
	testConn = conn
	code := m.Run()
	_ = testConn.Close(ctx)
	_ = container.Terminate(ctx)
	os.Exit(code)
}
