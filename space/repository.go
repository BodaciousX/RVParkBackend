// space/repository.go contains the repository interface for the space package.
package space

import (
	"database/sql"
)

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) List() ([]Space, error) {
	query := `
		SELECT 
			s.id,
			sec.name as section,
			s.status,
			s.tenant_id,
			s.reserved,
			s.payment_type,
			s.next_payment,
			s.tenant_notified,
			s.past_due_amount
		FROM spaces s
		JOIN sections sec ON s.section_id = sec.id
		ORDER BY sec.name, s.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spaces []Space
	for rows.Next() {
		var space Space
		var tenantID sql.NullString
		var paymentType sql.NullString
		var nextPayment sql.NullTime

		err := rows.Scan(
			&space.ID,
			&space.Section,
			&space.Status,
			&tenantID,
			&space.Reserved,
			&paymentType,
			&nextPayment,
			&space.TenantNotified,
			&space.PastDueAmount,
		)
		if err != nil {
			return nil, err
		}

		if tenantID.Valid {
			space.TenantID = &tenantID.String
		}
		if paymentType.Valid {
			space.PaymentType = paymentType.String
		}
		if nextPayment.Valid {
			space.NextPayment = nextPayment.Time
		}

		spaces = append(spaces, space)
	}

	return spaces, nil
}

func (r *sqlRepository) Get(id string) (*Space, error) {
	query := `
		SELECT 
			s.id,
			sec.name as section,
			s.status,
			s.tenant_id,
			s.reserved,
			s.payment_type,
			s.next_payment,
			s.tenant_notified,
			s.past_due_amount
		FROM spaces s
		JOIN sections sec ON s.section_id = sec.id
		WHERE s.id = $1
	`

	var space Space
	var tenantID sql.NullString
	var paymentType sql.NullString
	var nextPayment sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&space.ID,
		&space.Section,
		&space.Status,
		&tenantID,
		&space.Reserved,
		&paymentType,
		&nextPayment,
		&space.TenantNotified,
		&space.PastDueAmount,
	)
	if err != nil {
		return nil, err
	}

	if tenantID.Valid {
		space.TenantID = &tenantID.String
	}
	if paymentType.Valid {
		space.PaymentType = paymentType.String
	}
	if nextPayment.Valid {
		space.NextPayment = nextPayment.Time
	}

	return &space, nil
}

func (r *sqlRepository) Update(space Space) error {
	query := `
		UPDATE spaces SET
			status = $2,
			tenant_id = $3,
			reserved = $4,
			payment_type = $5,
			next_payment = $6,
			tenant_notified = $7,
			past_due_amount = $8
		WHERE id = $1
	`

	var tenantID interface{}
	if space.TenantID != nil {
		tenantID = *space.TenantID
	}

	var nextPayment interface{}
	if !space.NextPayment.IsZero() {
		nextPayment = space.NextPayment
	}

	_, err := r.db.Exec(
		query,
		space.ID,
		space.Status,
		tenantID,
		space.Reserved,
		space.PaymentType,
		nextPayment,
		space.TenantNotified,
		space.PastDueAmount,
	)
	return err
}
