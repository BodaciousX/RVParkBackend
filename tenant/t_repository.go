// tenant/t_repository.go
package tenant

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

func (r *sqlRepository) Create(tenant Tenant) error {
	query := `
        INSERT INTO tenants (
            id, name, move_in_date, space_id, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        )
    `

	now := time.Now()
	_, err := r.db.Exec(
		query,
		tenant.ID,
		tenant.Name,
		tenant.MoveInDate,
		tenant.SpaceID,
		now,
		now,
	)
	return err
}

func (r *sqlRepository) Get(id string) (*Tenant, error) {
	query := `
        SELECT 
            id,
            name,
            move_in_date,
            space_id,
            created_at,
            updated_at
        FROM tenants
        WHERE id = $1
    `

	var tenant Tenant
	err := r.db.QueryRow(query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.MoveInDate,
		&tenant.SpaceID,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r *sqlRepository) GetBySpace(spaceID string) (*Tenant, error) {
	query := `
        SELECT 
            id,
            name,
            move_in_date,
            space_id,
            created_at,
            updated_at
        FROM tenants
        WHERE space_id = $1
    `

	var tenant Tenant
	err := r.db.QueryRow(query, spaceID).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.MoveInDate,
		&tenant.SpaceID,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
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
            space_id = $3,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	_, err := r.db.Exec(
		query,
		tenant.ID,
		tenant.Name,
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
            move_in_date,
            space_id,
            created_at,
            updated_at
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
			&tenant.MoveInDate,
			&tenant.SpaceID,
			&tenant.CreatedAt,
			&tenant.UpdatedAt,
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
