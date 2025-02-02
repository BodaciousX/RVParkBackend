// user/u_repository.go contains the repository interface for the user package.
package user

import (
	"database/sql"
)

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) Create(user User) error {
	query := `
		INSERT INTO users (
			id,
			email,
			username,
			password_hash,
			role,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
	)
	return err
}

func (r *sqlRepository) Get(id string) (*User, error) {
	query := `
		SELECT 
			id,
			email,
			username,
			password_hash,
			role,
			created_at,
			last_login
		FROM users
		WHERE id = $1
	`

	var user User
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&lastLogin,
	)
	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time
	}

	return &user, nil
}

func (r *sqlRepository) GetByEmail(email string) (*User, error) {
	query := `
		SELECT 
			id,
			email,
			username,
			password_hash,
			role,
			created_at,
			last_login
		FROM users
		WHERE email = $1
	`

	var user User
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&lastLogin,
	)
	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time
	}

	return &user, nil
}

func (r *sqlRepository) Update(user User) error {
	query := `
		UPDATE users SET
			email = $2,
			username = $3,
			password_hash = $4,
			role = $5,
			last_login = $6
		WHERE id = $1
	`

	var lastLogin interface{}
	if !user.LastLogin.IsZero() {
		lastLogin = user.LastLogin
	}

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Role,
		lastLogin,
	)
	return err
}

func (r *sqlRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
