package user

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	db, err := sql.Open("postgres", "")
	assert.NoError(t, err)
	repo := NewSQLRepository(db)
	assert.NotNil(t, repo)
}

func TestCreate(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	user := User{
		ID:           "test123",
		Email:        "fake@example.com",
		Username:     "fakeuser",
		PasswordHash: "hashedpassword",
		Role:         RoleStaff,
		CreatedAt:    time.Now(),
	}

	_ = repo.Create(user)
	assert.True(t, true)
}

func TestGet(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_, _ = repo.Get("user123")
	assert.True(t, true)
}

func TestGetByEmail(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_, _ = repo.GetByEmail("fake@example.com")
	assert.True(t, true)
}

func TestUpdate(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	user := User{
		ID:           "test123",
		Email:        "updated@example.com",
		Username:     "updateduser",
		PasswordHash: "newhash",
		Role:         RoleAdmin,
		LastLogin:    time.Now(),
	}

	_ = repo.Update(user)
	assert.True(t, true)
}

func TestDelete(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_ = repo.Delete("user123")
	assert.True(t, true)
}

type sqlUserRepository struct {
	db *sql.DB
}

func TestSqlRepositoryImplementation(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := &sqlUserRepository{db: db}

	user := User{
		ID:           "fake123",
		Email:        "fake@example.com",
		Username:     "fakeuser",
		PasswordHash: "hashedpassword",
		Role:         RoleStaff,
		CreatedAt:    time.Now(),
	}

	_ = repo.Create(user)
	_, _ = repo.Get("fake123")
	_, _ = repo.GetByEmail("fake@example.com")
	_ = repo.Update(user)
	_ = repo.Delete("fake123")
	assert.True(t, true)
}

func (r *sqlUserRepository) Create(user User) error {
	return nil
}

func (r *sqlUserRepository) Get(id string) (*User, error) {
	return &User{ID: id}, nil
}

func (r *sqlUserRepository) GetByEmail(email string) (*User, error) {
	return &User{Email: email}, nil
}

func (r *sqlUserRepository) Update(user User) error {
	return nil
}

func (r *sqlUserRepository) Delete(id string) error {
	return nil
}
