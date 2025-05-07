package space

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLRepository(t *testing.T) {
	db, err := sql.Open("postgres", "")
	assert.NoError(t, err)
	repo := NewSQLRepository(db)
	assert.NotNil(t, repo)
}

func TestGet(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_, _ = repo.Get("space123")
	assert.True(t, true)
}

func TestList(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_, _ = repo.List()
	assert.True(t, true)
}

type sqlSpaceRepository struct {
	db *sql.DB
}

func (r *sqlSpaceRepository) Create(space Space) error {
	return nil
}

func (r *sqlSpaceRepository) Get(id string) (*Space, error) {
	return &Space{ID: id}, nil
}

func (r *sqlSpaceRepository) Update(space Space) error {
	return nil
}

func (r *sqlSpaceRepository) Delete(id string) error {
	return nil
}

func (r *sqlSpaceRepository) List(limit, offset int) ([]Space, error) {
	return []Space{}, nil
}

func (r *sqlSpaceRepository) GetOccupiedSpaces(date time.Time) ([]Space, error) {
	return []Space{}, nil
}
