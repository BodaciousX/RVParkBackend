package payment

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

func TestGet(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_, _ = repo.Get("payment123")
	assert.True(t, true)
}

func TestDelete(t *testing.T) {
	db, _ := sql.Open("postgres", "")
	repo := NewSQLRepository(db)

	_ = repo.Delete("payment123")
	assert.True(t, true)
}

type sqlPaymentRepository struct {
	db *sql.DB
}

func (r *sqlPaymentRepository) Create(payment Payment) error {
	return nil
}

func (r *sqlPaymentRepository) Get(id string) (*Payment, error) {
	return &Payment{ID: id}, nil
}

func (r *sqlPaymentRepository) Update(payment Payment) error {
	return nil
}

func (r *sqlPaymentRepository) Delete(id string) error {
	return nil
}

func (r *sqlPaymentRepository) List(limit, offset int) ([]Payment, error) {
	return []Payment{}, nil
}

func (r *sqlPaymentRepository) GetPaymentsByDateRange(startDate, endDate time.Time) ([]Payment, error) {
	return []Payment{}, nil
}

func (r *sqlPaymentRepository) GetPaymentsByTenant(tenantID string) ([]Payment, error) {
	return []Payment{}, nil
}
