package conn

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	connStr := "postgres://a1111:dkfl26052010@localhost:5432/postgres?sslmode=disable"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Println("failed to connect postgres")
	}
	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, err
	}
	query := `
	CREATE TABLE IF NOT EXISTS users(
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(70) NOT NULL,
	email VARCHAR(50) NOT NULL,
	balance NUMERIC(19, 2) NOT NULL DEFAULT 0
	)
	`
	if _, err := conn.Exec(ctx, query); err != nil {
		fmt.Println("failed to create table")
		return nil, err
	}
	return conn, nil
}
