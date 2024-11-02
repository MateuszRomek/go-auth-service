package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mateuszromek/auth/internal/auth"
)

type Session struct {
	Id        string `db:"id" json:"id"`
	UserId    string `db:"user_id" json:"user_id"`
	ExpiresAt int64  `db:"expires_at" json:"expires_at"`
}

type SessionStore struct {
	db *sqlx.DB
}

func (s *SessionStore) CreateSession(ctx context.Context, userId string) (*Session, string, error) {
	token, err := auth.GenerateSessionToken()
	fmt.Println(userId, token)
	if err != nil {
		return nil, "", err
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
		return nil, "", err
	}

	defer rows.Close()

	return &session, token, nil
}

func (s *SessionStore) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	sessionId := auth.CreateSessionId(token)

	query := `
		SELECT * FROM session WHERE id = $1 LIMIT 1
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(ctxWithTimeout, query, sessionId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("failed to get session")
	}

	var session Session
	err = rows.StructScan(&session)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionStore) ValidateSessionByToken(ctx context.Context, token string) (*Session, string, error) {
	sessionId := auth.CreateSessionId(token)

	query := `
	SELECT id, user_id, expires_at
	FROM session
	WHERE id = $1
	`

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(ctxWithTimeout, query, sessionId)
	if err != nil {
		return nil, "", err
	}

	if !rows.Next() {
		return nil, "", errors.New("no session for provided token")
	}

	var session Session
	err = rows.StructScan(&session)

	if err != nil {
		return nil, "", err
	}

	now := time.Now().UnixMilli()

	if now >= session.ExpiresAt {
		deleteQuery := `
		 DELETE FROM session WHERE id = $1
		`

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_, err := s.db.ExecContext(ctxWithTimeout, deleteQuery, session.Id)

		if err != nil {
			return nil, "", errors.New("failed to remove expired session")
		}

		return nil, "", errors.New("token expired")
	}

	newExpiresAt := time.Now().Add(7 * 24 * time.Hour).UnixMilli()
	updateQuery := `
	  UPDATE session SET expires_at = $1 WHERE id = $2 
		RETURNING id, user_id, expires_at
	`

	ctxWithTimeout, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err = s.db.QueryxContext(ctxWithTimeout, updateQuery, newExpiresAt, session.Id)

	if err != nil {
		return nil, "", errors.New("failed to extend query")
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, "", errors.New("no session found")
	}

	var updatedSession Session
	err = rows.StructScan(&updatedSession)
	if err != nil {
		return nil, "", errors.New("failed to scan updated session")
	}

	return &updatedSession, token, nil
}

func (s *SessionStore) RemoveAllUserSessions(ctx context.Context, userId string) error {
	query := `
	DELETE FROM session
	WHERE user_id = $1
	`
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctxWithTimeout, query, userId)
	if err != nil {
		return errors.New("failed to remove user sessions")
	}

	return nil
}
