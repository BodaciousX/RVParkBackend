// tenant/repository.go contains the repository interface for the tenant package.
package tenant

import (
	"database/sql"
)

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) Create(tenant Tenant) error {
	query := `
		INSERT INTO tenants (
			id, name, phone, move_in_date, space_id
		) VALUES (
			$1, $2, $3, $4, $5
		)
	`

	_, err := r.db.Exec(
		query,
		tenant.ID,
		tenant.Name,
		tenant.Phone,
		tenant.MoveInDate,
		tenant.SpaceID,
	)
	return err
}

func (r *sqlRepository) Get(id string) (*Tenant, error) {
	query := `
		SELECT 
			id,
			name,
			phone,
			move_in_date,
			space_id
		FROM tenants
		WHERE id = $1
	`

	var tenant Tenant
	err := r.db.QueryRow(query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Phone,
		&tenant.MoveInDate,
		&tenant.SpaceID,
	)
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *sqlRepository) Update(tenant Tenant) error {
	query := `
		UPDATE tenants SET
			name = $2,
			phone = $3,
			space_id = $4
		WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		tenant.ID,
		tenant.Name,
		tenant.Phone,
		tenant.SpaceID,
	)
	return err
}

func (r *sqlRepository) Delete(id string) error {
	query := `DELETE FROM tenants WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *sqlRepository) List() ([]Tenant, error) {
	query := `
		SELECT 
			id,
			name,
			phone,
			move_in_date,
			space_id
		FROM tenants
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []Tenant
	for rows.Next() {
		var tenant Tenant
		err := rows.Scan(
			&tenant.ID,
			&tenant.Name,
			&tenant.Phone,
			&tenant.MoveInDate,
			&tenant.SpaceID,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

func (r *sqlRepository) ListPayments(tenantID string) ([]Payment, error) {
	query := `
		SELECT 
			id,
			tenant_id,
			amount,
			due_date,
			paid_date,
			payment_type,
			status
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
			&payment.Amount,
			&payment.DueDate,
			&paidDate,
			&payment.PaymentType,
			&payment.Status,
		)
		if err != nil {
			return nil, err
		}

		if paidDate.Valid {
			payment.PaidDate = &paidDate.Time
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *sqlRepository) CreatePayment(payment Payment) error {
	query := `
		INSERT INTO payments (
			id,
			tenant_id,
			amount,
			due_date,
			paid_date,
			payment_type,
			status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	_, err := r.db.Exec(
		query,
		payment.ID,
		payment.TenantID,
		payment.Amount,
		payment.DueDate,
		payment.PaidDate,
		payment.PaymentType,
		payment.Status,
	)
	return err
}

func (r *sqlRepository) UpdatePayment(payment Payment) error {
	query := `
		UPDATE payments SET
			amount = $2,
			paid_date = $3,
			status = $4
		WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		payment.ID,
		payment.Amount,
		payment.PaidDate,
		payment.Status,
	)
	return err
}
