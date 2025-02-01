// space/service.go
package space

import (
	"fmt"

	"github.com/BodaciousX/RVParkBackend/tenant"
)

type service struct {
	repo          Repository
	tenantService tenant.Service
}

func NewService(repo Repository, tenantService tenant.Service) Service {
	return &service{
		repo:          repo,
		tenantService: tenantService,
	}
}

func (s *service) ListSpaces() (map[string][]Space, error) {
	spaces, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// Group spaces by section
	grouped := make(map[string][]Space)
	for _, space := range spaces {
		grouped[space.Section] = append(grouped[space.Section], space)
	}
	return grouped, nil
}

func (s *service) GetSpace(id string) (*Space, error) {
	space, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return space, nil
}

func (s *service) GetVacantSpaces() ([]Space, error) {
	spaces, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	var vacant []Space
	for _, space := range spaces {
		if space.TenantID == nil && !space.Reserved {
			vacant = append(vacant, space)
		}
	}
	return vacant, nil
}

func (s *service) ReserveSpace(spaceID string) error {
	space, err := s.repo.Get(spaceID)
	if err != nil {
		return err
	}

	// Can only reserve vacant spaces
	if space.Status != StatusVacant {
		return fmt.Errorf("space %s is not vacant", spaceID)
	}

	space.Reserved = true
	space.Status = StatusReserved
	return s.repo.Update(*space)
}

func (s *service) UnreserveSpace(spaceID string) error {
	space, err := s.repo.Get(spaceID)
	if err != nil {
		return err
	}

	// Can only unreserve reserved spaces
	if space.Status != StatusReserved {
		return fmt.Errorf("space %s is not reserved", spaceID)
	}

	space.Reserved = false
	space.Status = StatusVacant
	return s.repo.Update(*space)
}

func (s *service) MoveIn(spaceID string, tenantID string) error {
	space, err := s.repo.Get(spaceID)
	if err != nil {
		return err
	}

	// Can only move in to vacant or reserved spaces
	if space.Status != StatusVacant && space.Status != StatusReserved {
		return fmt.Errorf("space %s is not available", spaceID)
	}

	space.TenantID = &tenantID
	space.Status = StatusOccupied
	space.Reserved = false

	return s.repo.Update(*space)
}

func (s *service) MoveOut(spaceID string) error {
	space, err := s.repo.Get(spaceID)
	if err != nil {
		return err
	}

	// Can only move out from occupied spaces
	if space.Status != StatusOccupied {
		return fmt.Errorf("space %s is not occupied", spaceID)
	}

	space.TenantID = nil
	space.Status = StatusVacant
	space.Reserved = false

	return s.repo.Update(*space)
}

func (s *service) UpdateSpace(space Space) error {
	// Validate status
	switch space.Status {
	case StatusOccupied, StatusVacant, StatusReserved:
		// Valid status
	default:
		return fmt.Errorf("invalid status: %s", space.Status)
	}

	// Validate state consistency
	if space.Reserved && space.Status != StatusReserved {
		return fmt.Errorf("reserved spaces must have Reserved status")
	}

	if space.TenantID != nil && space.Status != StatusOccupied {
		return fmt.Errorf("spaces with tenants must have Occupied status")
	}

	if space.Status == StatusOccupied && space.TenantID == nil {
		return fmt.Errorf("occupied spaces must have a tenant")
	}

	return s.repo.Update(space)
}
