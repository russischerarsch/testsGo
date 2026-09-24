package repo

import (
	"context"
	"strconv"
	"tgtest/domain"
	"tgtest/serv"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func CreateRepo(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, eventID string, user *domain.User) (string, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	tx.Begin(ctx)
	defer tx.Rollback(ctx)
	query := `
	INSERT INTO users (name, email)
	VALUES($1, $2)
	RETURNING id
	`
	if err := r.db.QueryRow(ctx, query, user.Name, user.Email).Scan(&user.Id); err != nil {
		return "0", err
	}
	var event = &serv.UserCreatedEvent{
		UserID:    strconv.Itoa(user.Id),
		EventID:   eventID,
		Action:    "creation",
		CreatedAt: time.Now().UTC(),
	}
	query = `
	INSERT INTO events (user_id, event_id, action, created_at)
	`
	r.db.Exec(ctx)
	return strconv.Itoa(user.Id), nil
}
