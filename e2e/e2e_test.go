package e2e

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"tgtest/http"
	"tgtest/repo"
	"tgtest/serv"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ctx = context.Background()
var testConn *pgx.Conn
var serverURL string

func TestMain(t *testing.M) {
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
		log.Fatalf("failed to set up testcontainers, %v", err)
	}
	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	url := "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatalf("failed to open connection %v", err)
	}
	testConn = conn
	m, err := migrate.New("file://../migrations", url)
	if err != nil {
		log.Fatalf("failed to find migration file %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to set migrations %v", err)
	}
	m.Close()
	repository := repo.CreateRepo(conn)
	service, err := serv.CreateServ(repository)
	if err != nil {
		log.Fatalf("failed to create service, %v", err)
	}
	handler := http.CreateHandler(service)
	router := gin.Default()
	router.POST("/users", handler.CreateUser)
	ts := httptest.NewServer(router)
	serverURL = ts.URL
	code := t.Run()
	ts.Close()
	_ = conn.Close(ctx)
	_ = pgContainer.Terminate(ctx)
	os.Exit(code)
}
