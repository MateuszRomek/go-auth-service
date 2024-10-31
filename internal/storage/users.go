package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mateuszromek/auth/internal/auth"
)

type User struct {
	Id            string `db:"id"`
	Email         string `db:"email"`
	PasswordHash  string `db:"password_hash"`
	EmailVerified bool   `db:"email_verified"`
}

type CreateUserPayload struct {
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

type UserStore struct {
	db *sqlx.DB
}

func handleSqlError(err error) string {
	var pqError *pq.Error

	if errors.As(err, &pqError) {
		switch pqError.Code {
		case "23505":
			if pqError.Constraint == "user_account_email_key" {
				return "Email already exist"
			}

			return "User already exist"
		}
	}

	return "Internal server error"
}

func (s *UserStore) CreateUser(ctx context.Context, email, password string) (*User, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO user_account (email, password_hash)
		VALUES (:email, :password_hash)
		RETURNING id, email, password_hash, email_verified
	`

	createUserPayload := CreateUserPayload{
		Email:        email,
		PasswordHash: hashedPassword,
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.NamedQueryContext(ctxWithTimeout, query, createUserPayload)
	if err != nil {
		return nil, errors.New(handleSqlError(err))
	}
	defer rows.Close()

	var user User
	if !rows.Next() {
		return nil, errors.New("failed to create user")
	}

	err = rows.StructScan(&user)
	if err != nil {
		return nil, err
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

	if !rows.Next() {
		return nil, errors.New("failed to create user")
	}
	var user User

	err = rows.StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
