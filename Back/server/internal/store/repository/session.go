package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSessionAlreadyRotated = errors.New("session was already rotated")

type SessionRepository struct {
	db *pgxpool.Pool
}

type Session struct {
	ID           int64
	UserID       int
	Username     string
	FamilyID     string
	IsRevoked    bool
	ExpiresAt    time.Time
	RevokedAt    *time.Time
	ReplacedByID *int64
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, userID int, tokenHash, familyID string, expiresAt time.Time, ipAddress, userAgent string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO sessions (user_id, token_hash, family_id, expires_at, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, tokenHash, familyID, expiresAt, ipAddress, userAgent,
	)
	return err
}

func (r *SessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*Session, error) {
	var session Session
	err := r.db.QueryRow(ctx,
		`SELECT s.id, s.user_id, u.username, s.family_id, s.is_revoked, s.expires_at,
		        s.revoked_at, s.replaced_by_id
		 FROM sessions s
		 JOIN "user" u ON u.id = s.user_id
		 WHERE s.token_hash = $1`,
		tokenHash,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.Username,
		&session.FamilyID,
		&session.IsRevoked,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.ReplacedByID,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Rotate revokes oldSessionID and inserts a new refresh token in the same
// family. Concurrent refreshes race on the conditional UPDATE: only one wins,
// the rest receive ErrSessionAlreadyRotated.
func (r *SessionRepository) Rotate(ctx context.Context, oldSessionID int64, newTokenHash string, expiresAt time.Time, ipAddress, userAgent string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var newSessionID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO sessions (user_id, token_hash, family_id, expires_at, ip_address, user_agent)
		 SELECT user_id, $2, family_id, $3, $4, $5
		 FROM sessions
		 WHERE id = $1
		 RETURNING id`,
		oldSessionID, newTokenHash, expiresAt, ipAddress, userAgent,
	).Scan(&newSessionID)
	if err != nil {
		return err
	}

	result, err := tx.Exec(ctx,
		`UPDATE sessions
		 SET is_revoked = TRUE, revoked_at = CURRENT_TIMESTAMP, replaced_by_id = $2
		 WHERE id = $1 AND is_revoked = FALSE`,
		oldSessionID, newSessionID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrSessionAlreadyRotated
	}

	return tx.Commit(ctx)
}

func (r *SessionRepository) RevokeFamily(ctx context.Context, familyID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE sessions
		 SET is_revoked = TRUE, revoked_at = CURRENT_TIMESTAMP
		 WHERE family_id = $1 AND is_revoked = FALSE`,
		familyID,
	)
	return err
}

func (r *SessionRepository) Revoke(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE sessions
		 SET is_revoked = TRUE, revoked_at = CURRENT_TIMESTAMP
		 WHERE token_hash = $1 AND is_revoked = FALSE`,
		tokenHash,
	)
	return err
}

func (r *SessionRepository) RevokeAllForUser(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE sessions
		 SET is_revoked = TRUE, revoked_at = CURRENT_TIMESTAMP
		 WHERE user_id = $1 AND is_revoked = FALSE`,
		userID,
	)
	return err
}

func (r *SessionRepository) DeleteExpired(ctx context.Context, now time.Time) (int64, error) {
	result, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE expires_at < $1`, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
