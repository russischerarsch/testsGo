package repo

import (
	"context"
	"strconv"
	"tgtest/domain"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func CreateRepo(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	query := `
	INSERT INTO users (name, email)
	VALUES($1, $2)
	RETURNING id
	`
	if err := r.db.QueryRow(ctx, query, user.Name, user.Email).Scan(&user.Id); err != nil {
		return "0", err
	}

	return strconv.Itoa(user.Id), nil
}
