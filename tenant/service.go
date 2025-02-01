// tenant/service.go
package tenant

import (
	"fmt"
	"time"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTenant(tenant Tenant) error {
	// Validate tenant data
	if tenant.Name == "" {
		return fmt.Errorf("tenant name is required")
	}
	if tenant.SpaceID == "" {
		return fmt.Errorf("space ID is required")
	}
	if tenant.MoveInDate.IsZero() {
		tenant.MoveInDate = time.Now()
	}

	// Check if space already has a tenant
	existingTenant, err := s.repo.GetBySpace(tenant.SpaceID)
	if err == nil && existingTenant != nil {
		return fmt.Errorf("space %s is already occupied", tenant.SpaceID)
	}

	return s.repo.Create(tenant)
}

func (s *service) GetTenant(id string) (*Tenant, error) {
	return s.repo.Get(id)
}

func (s *service) GetTenantBySpace(spaceID string) (*Tenant, error) {
	return s.repo.GetBySpace(spaceID)
}

func (s *service) UpdateTenant(tenant Tenant) error {
	// Validate tenant exists
	existing, err := s.repo.Get(tenant.ID)
	if err != nil {
		return fmt.Errorf("tenant not found: %v", err)
	}

	// If space is changing, check if new space is available
	if tenant.SpaceID != existing.SpaceID {
		existingTenant, err := s.repo.GetBySpace(tenant.SpaceID)
		if err == nil && existingTenant != nil {
			return fmt.Errorf("space %s is already occupied", tenant.SpaceID)
		}
	}

	// Preserve creation time and move-in date
	tenant.CreatedAt = existing.CreatedAt
	tenant.MoveInDate = existing.MoveInDate

	return s.repo.Update(tenant)
}

func (s *service) DeleteTenant(id string) error {
	// Verify tenant exists before deletion
	if _, err := s.repo.Get(id); err != nil {
		return fmt.Errorf("tenant not found: %v", err)
	}

	return s.repo.Delete(id)
}

func (s *service) ListTenants() ([]Tenant, error) {
	return s.repo.List()
}
