// payment/p_repository.go
package payment

import (
	"database/sql"
	"time"
)

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) Create(payment Payment) error {
	query := `
        INSERT INTO payments (
            id, tenant_id, amount_due, due_date, paid_date, 
            next_payment_date, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
    `

	now := time.Now()
	_, err := r.db.Exec(
		query,
		payment.ID,
		payment.TenantID,
		payment.AmountDue,
		payment.DueDate,
		payment.PaidDate,
		payment.NextPaymentDate,
		now,
	)
	return err
}

func (r *sqlRepository) Get(id string) (*Payment, error) {
	query := `
        SELECT 
            id, tenant_id, amount_due, due_date, paid_date,
            next_payment_date, created_at, updated_at
        FROM payments
        WHERE id = $1
    `

	var payment Payment
	var paidDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.TenantID,
		&payment.AmountDue,
		&payment.DueDate,
		&paidDate,
		&payment.NextPaymentDate,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if paidDate.Valid {
		payment.PaidDate = &paidDate.Time
	}

	return &payment, nil
}

func (r *sqlRepository) Update(payment Payment) error {
	query := `
        UPDATE payments SET
            amount_due = $2,
            due_date = $3,
            paid_date = $4,
            next_payment_date = $5,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	_, err := r.db.Exec(
		query,
		payment.ID,
		payment.AmountDue,
		payment.DueDate,
		payment.PaidDate,
		payment.NextPaymentDate,
	)
	return err
}

func (r *sqlRepository) Delete(id string) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *sqlRepository) ListByTenant(tenantID string) ([]Payment, error) {
	query := `
        SELECT 
            id, tenant_id, amount_due, due_date, paid_date,
            next_payment_date, created_at, updated_at
        FROM payments
        WHERE tenant_id = $1
        ORDER BY due_date DESC
    `

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		var paidDate sql.NullTime

		err := rows.Scan(
			&payment.ID,
			&payment.TenantID,
			&payment.AmountDue,
			&payment.DueDate,
			&paidDate,
			&payment.NextPaymentDate,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if paidDate.Valid {
			payment.PaidDate = &paidDate.Time
		}

		payments = append(payments, payment)
	}

	return payments, rows.Err()
}

func (r *sqlRepository) ListByDateRange(start, end time.Time) ([]Payment, error) {
	query := `
        SELECT 
            id, tenant_id, amount_due, due_date, paid_date,
            next_payment_date, created_at, updated_at
        FROM payments
        WHERE due_date BETWEEN $1 AND $2
        ORDER BY due_date DESC
    `

	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		var paidDate sql.NullTime

		err := rows.Scan(
			&payment.ID,
			&payment.TenantID,
			&payment.AmountDue,
			&payment.DueDate,
			&paidDate,
			&payment.NextPaymentDate,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if paidDate.Valid {
			payment.PaidDate = &paidDate.Time
		}

		payments = append(payments, payment)
	}

	return payments, rows.Err()
}

func (r *sqlRepository) GetLatestByTenant(tenantID string) (*Payment, error) {
	query := `
        SELECT 
            id, tenant_id, amount_due, due_date, paid_date,
            next_payment_date, created_at, updated_at
        FROM payments
        WHERE tenant_id = $1
        ORDER BY due_date DESC
        LIMIT 1
    `

	var payment Payment
	var paidDate sql.NullTime

	err := r.db.QueryRow(query, tenantID).Scan(
		&payment.ID,
		&payment.TenantID,
		&payment.AmountDue,
		&payment.DueDate,
		&paidDate,
		&payment.NextPaymentDate,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if paidDate.Valid {
		payment.PaidDate = &paidDate.Time
	}

	return &payment, nil
}
