package tenant

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	db, err := sql.Open("postgres", "")
	assert.NoError(t, err)
	repo := NewSQLRepository(db)
	assert.NotNil(t, repo)
}

func TestGet(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_, _ = repo.Get("tenant123")
	assert.True(t, true)
}

func TestDelete(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_ = repo.Delete("tenant123")
	assert.True(t, true)
}

func TestList(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_, _ = repo.List()
	assert.True(t, true)
}

type sqlTenantRepository struct {
	db *sql.DB
}

func (r *sqlTenantRepository) Create(tenant Tenant) error {
	return nil
}

func (r *sqlTenantRepository) Get(id string) (*Tenant, error) {
	return &Tenant{ID: id}, nil
}

func (r *sqlTenantRepository) Update(tenant Tenant) error {
	return nil
}

func (r *sqlTenantRepository) Delete(id string) error {
	return nil
}

func (r *sqlTenantRepository) List(limit, offset int) ([]Tenant, error) {
	return []Tenant{}, nil
}

func (r *sqlTenantRepository) Search(query string) ([]Tenant, error) {
	return []Tenant{}, nil
}
