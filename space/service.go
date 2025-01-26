// space/service.go contains the implementation of the Service interface.
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

	// If space has tenant, get payment status
	if space.TenantID != nil {
		status, pastDue, err := s.tenantService.GetPaymentStatus(*space.TenantID)
		if err != nil {
			return nil, err
		}
		space.Status = "Occupied (" + status + ")"
		space.PastDueAmount = pastDue
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

	space.Reserved = true
	return s.repo.Update(*space)
}

func (s *service) UnreserveSpace(spaceID string) error {
	space, err := s.repo.Get(spaceID)
	if err != nil {
		return err
	}

	space.Reserved = false
	return s.repo.Update(*space)
}

func (s *service) MoveIn(spaceID string, tenantID string) error {
	space, err := s.repo.Get(spaceID)
	if err != nil {
		return err
	}

	space.TenantID = &tenantID
	space.Status = "Occupied (Paid)"
	space.Reserved = false
	return s.repo.Update(*space)
}

func (s *service) MoveOut(spaceID string) error {
	space, err := s.repo.Get(spaceID)
	if err != nil {
		return err
	}

	space.TenantID = nil
	space.Status = "Vacant"
	space.PaymentType = ""
	return s.repo.Update(*space)
}

func (s *service) UpdateSpace(space Space) error {
	// Validate status
	switch space.Status {
	case "Occupied (Paid)", "Occupied (Payment Due)", "Occupied (Overdue)", "Vacant", "Reserved":
		// Valid status
	default:
		return fmt.Errorf("invalid status: %s", space.Status)
	}

	return s.repo.Update(space)
}
