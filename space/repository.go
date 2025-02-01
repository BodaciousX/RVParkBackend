// space/repository.go
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
            s.reserved
        FROM spaces s
        JOIN sections sec ON s.section_id = sec.id
        ORDER BY 
            sec.name,
            SUBSTRING(s.id FROM '^[A-Za-z]+'),
            CAST(SUBSTRING(s.id FROM '[0-9]+') AS INTEGER)
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

		err := rows.Scan(
			&space.ID,
			&space.Section,
			&space.Status,
			&tenantID,
			&space.Reserved,
		)
		if err != nil {
			return nil, err
		}

		if tenantID.Valid {
			space.TenantID = &tenantID.String
		}

		spaces = append(spaces, space)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
            s.reserved
        FROM spaces s
        JOIN sections sec ON s.section_id = sec.id
        WHERE s.id = $1
    `

	var space Space
	var tenantID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&space.ID,
		&space.Section,
		&space.Status,
		&tenantID,
		&space.Reserved,
	)
	if err != nil {
		return nil, err
	}

	if tenantID.Valid {
		space.TenantID = &tenantID.String
	}

	return &space, nil
}

func (r *sqlRepository) Update(space Space) error {
	query := `
        UPDATE spaces SET
            status = $2,
            tenant_id = $3,
            reserved = $4,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `

	var tenantID interface{}
	if space.TenantID != nil {
		tenantID = *space.TenantID
	}

	_, err := r.db.Exec(
		query,
		space.ID,
		space.Status,
		tenantID,
		space.Reserved,
	)
	return err
}
