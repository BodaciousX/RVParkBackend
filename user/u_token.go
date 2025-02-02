// user/u_token.go contains the implementation of the user token.
package user

import (
	"database/sql"
	"errors"
	"time"
)

type Token struct {
	TokenHash string    `json:"tokenHash"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	Revoked   bool      `json:"revoked"`
}

type TokenRepository interface {
	CreateToken(token Token) error
	GetToken(tokenHash string) (*Token, error)
	RevokeToken(tokenHash string) error
	RevokeAllUserTokens(userID string) error
	CleanExpiredTokens() error
}

type sqlTokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &sqlTokenRepository{db: db}
}

// CreateToken stores a new token in the database
func (r *sqlTokenRepository) CreateToken(token Token) error {
	query := `
		INSERT INTO tokens (token_hash, user_id, expires_at, created_at, revoked)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(
		query,
		token.TokenHash,
		token.UserID,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	)
	return err
}

// GetToken retrieves a token by its hash
func (r *sqlTokenRepository) GetToken(tokenHash string) (*Token, error) {
	query := `
		SELECT token_hash, user_id, expires_at, created_at, revoked
		FROM tokens
		WHERE token_hash = $1
	`

	token := &Token{}
	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.TokenHash,
		&token.UserID,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.Revoked,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("token not found")
	}
	if err != nil {
		return nil, err
	}

	return token, nil
}

// RevokeToken marks a specific token as revoked
func (r *sqlTokenRepository) RevokeToken(tokenHash string) error {
	query := `
		UPDATE tokens
		SET revoked = true
		WHERE token_hash = $1
	`
	result, err := r.db.Exec(query, tokenHash)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("token not found")
	}

	return nil
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (r *sqlTokenRepository) RevokeAllUserTokens(userID string) error {
	query := `
		UPDATE tokens
		SET revoked = true
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	return err
}

// CleanExpiredTokens removes all expired and revoked tokens from the database
func (r *sqlTokenRepository) CleanExpiredTokens() error {
	query := `
		DELETE FROM tokens
		WHERE expires_at < $1
		OR revoked = true
	`
	_, err := r.db.Exec(query, time.Now())
	return err
}
