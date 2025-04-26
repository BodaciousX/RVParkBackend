// tenant/service.go
package tenant

import (
	"fmt"
	"time"

	"github.com/BodaciousX/RVParkBackend/space"
)

type service struct {
	repo         Repository
	spaceService space.Service
}

func NewService(repo Repository, spaceService space.Service) Service {
	return &service{
		repo:         repo,
		spaceService: spaceService,
	}
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

	// Verify the space exists and is available
	space, err := s.spaceService.GetSpace(tenant.SpaceID)
	if err != nil {
		return fmt.Errorf("invalid space ID: %v", err)
	}

	// Check if space is available
	if space.Status != "Vacant" && space.Status != "Reserved" {
		return fmt.Errorf("space %s is not available", tenant.SpaceID)
	}

	// Create the tenant record first
	if err := s.repo.Create(tenant); err != nil {
		return err
	}

	// Then update the space to assign the tenant
	if err := s.spaceService.MoveIn(tenant.SpaceID, tenant.ID); err != nil {
		// If space update fails, attempt to rollback tenant creation
		s.repo.Delete(tenant.ID)
		return fmt.Errorf("failed to update space: %v", err)
	}

	return nil
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

	// If space is changing, update space assignments
	if tenant.SpaceID != existing.SpaceID {
		// Get the new space
		newSpace, err := s.spaceService.GetSpace(tenant.SpaceID)
		if err != nil {
			return fmt.Errorf("invalid space ID: %v", err)
		}

		// Check if new space is available
		if newSpace.Status != "Vacant" && newSpace.Status != "Reserved" {
			return fmt.Errorf("space %s is not available", tenant.SpaceID)
		}

		// Move out from the old space
		if err := s.spaceService.MoveOut(existing.SpaceID); err != nil {
			return fmt.Errorf("failed to vacate current space: %v", err)
		}

		// Move into the new space
		if err := s.spaceService.MoveIn(tenant.SpaceID, tenant.ID); err != nil {
			// Try to restore the old space assignment if this fails
			s.spaceService.MoveIn(existing.SpaceID, tenant.ID)
			return fmt.Errorf("failed to assign new space: %v", err)
		}
	}

	// Preserve creation time and move-in date
	tenant.CreatedAt = existing.CreatedAt
	tenant.MoveInDate = existing.MoveInDate

	return s.repo.Update(tenant)
}

func (s *service) DeleteTenant(id string) error {
	// Get the tenant to find their space
	tenant, err := s.repo.Get(id)
	if err != nil {
		return fmt.Errorf("tenant not found: %v", err)
	}

	// Remove tenant from space
	if err := s.spaceService.MoveOut(tenant.SpaceID); err != nil {
		return fmt.Errorf("failed to vacate space: %v", err)
	}

	// Then delete the tenant record
	return s.repo.Delete(id)
}

func (s *service) ListTenants() ([]Tenant, error) {
	return s.repo.List()
}
