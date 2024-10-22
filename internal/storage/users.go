package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mateuszromek/auth/internal/auth"
)

type User struct {
	Id            string `db:"id"`
	Email         string `db:"email"`
	Username      string `db:"username"`
	PasswordHash  string `db:"password_hash"`
	Salt          string `db:"salt"`
	EmailVerified bool   `db:"email_verified"`
}

type CreateUserPayload struct {
	Email        string `db:"email"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Salt         string `db:"salt"`
}

type UserStore struct {
	db *sqlx.DB
}

func (s *UserStore) CreateUser(ctx context.Context, email, username, password string) (*User, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO user_account (email, username, password_hash, salt)
		VALUES (:email, :username, :password_hash, :salt)
	`

	createUserPayload := CreateUserPayload{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword[1],
		Salt:         hashedPassword[0],
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.NamedQueryContext(ctxWithTimeout, query, createUserPayload)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user User
	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	return &user, nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT * FROM user_account WHERE email = $1`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(ctxWithTimeout, query, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var user User
	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

	}
	return &user, nil
}
