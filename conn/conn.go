package conn

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func CreateConnection(connStr string, ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("подключение к PostgreSQL: %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, err
	}
	return conn, nil
}
