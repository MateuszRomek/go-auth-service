package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	Sessions interface {
		CreateSession(ctx context.Context, userId string) (*Session, error)
		GetUserSession(ctx context.Context, userId string) (*Session, error)
		GetSessionByToken(ctx context.Context, token string) (*Session, error)
	}
	Users interface {
		CreateUser(ctx context.Context, email, username, password string) (*User, error)
		GetUserByEmail(ctx context.Context, email string) (*User, error)
	}
}

func NewStorage(db *sqlx.DB) Storage {
	return Storage{
		Sessions: &SessionStore{db},
		Users:    &UserStore{db},
	}
}
