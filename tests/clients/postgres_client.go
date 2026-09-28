package clients

import (
	"context"
	"fmt"
	"os"
	"tgtest/conn"

	"github.com/jackc/pgx/v5"
)

type PostgresClient struct {
	db *pgx.Conn
}

func CreatePostgresClient(dsn string, ctx context.Context) (*PostgresClient, error) {

	connStr := os.Getenv("POSTGRES_CONN")
	conn, err := conn.CreateConnection(connStr, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres, %w", err)
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to create connection, %w", err)
	}
	return &PostgresClient{db: conn}, nil
}
func (c *PostgresClient) Close(ctx context.Context) {
	c.db.Close(ctx)
}
func (c *PostgresClient) GetUserID(ctx context.Context, name string) (string, error) {
	query := `
	SELECT id FROM users 
	WHERE name = $1
	`
	var id string
	if err := c.db.QueryRow(ctx, query, name).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
func (c *PostgresClient) DeleteUser(ctx context.Context, id string) (bool, error) {
	query := `
	DELETE FROM users 
	WHERE id = $1
	`
	row, err := c.db.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}
	if row.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}
