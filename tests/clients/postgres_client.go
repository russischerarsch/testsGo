package clients

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresClient struct {
	db *sql.DB
}

func CreatePostgresClient(dsn string) (*PostgresClient, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres, %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to create connection, %w", err)
	}
	return &PostgresClient{db: conn}, nil
}
func (c *PostgresClient) Close() {
	c.db.Close()
}
func (c *PostgresClient) GetUserID(ctx context.Context, name string) (string, error) {
	query := `
	SELECT id FROM users 
	WHERE name = $1
	`
	var id string
	if err := c.db.QueryRowContext(ctx, query, name).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
