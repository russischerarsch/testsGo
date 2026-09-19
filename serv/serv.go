package serv

import (
	"context"
	"tgtest/domain"
)

type RepoInterface interface {
	CreateUser(ctx context.Context, user *domain.User) (int, error)
}
