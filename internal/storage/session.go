package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mateuszromek/auth/internal/auth"
)

type Session struct {
	Id        string `db:"id"`
	UserId    string `db:"user_id"`
	ExpiresAt int64  `db:"expires_at"`
}

type SessionStore struct {
	db *sqlx.DB
}

func (s *SessionStore) CreateSession(ctx context.Context, userId string) (*Session, error) {
	token, err := auth.GenerateSessionToken()
	if err != nil {
		return nil, err
	}

	sessionId := auth.CreateSessionId(token)

	query := `
		INSERT INTO session (id, user_id, expires_at)
		VALUES (:id, :user_id, :expires_at)
	`

	session := Session{
		Id:        sessionId,
		UserId:    userId,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).UnixMilli(),
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.NamedQueryContext(ctxWithTimeout, query, session)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return &session, nil
}

func (s *SessionStore) GetUserSession(ctx context.Context, userId string) (*Session, error) {
	query := `
		SELECT * FROM session WHERE user_id = $1
 		`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(ctxWithTimeout, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var session Session
	if rows.Next() {
		err := rows.StructScan(&session)
		if err != nil {
			return nil, err
		}
	}

	return &session, nil
}

func (s *SessionStore) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	sessionId := auth.CreateSessionId(token)

	query := `
		SELECT * FROM session WHERE id = $1
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(ctxWithTimeout, query, sessionId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var session Session
	if rows.Next() {
		err := rows.StructScan(&session)
		if err != nil {
			return nil, err
		}
	}

	return &session, nil
}
