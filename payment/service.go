// payment/service.go
package payment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePayment(payment Payment) error {
	// Validate payment data
	if payment.TenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}
	if payment.SpaceID == "" {
		return fmt.Errorf("space ID is required")
	}
	if payment.AmountDue <= 0 {
		return fmt.Errorf("amount due must be greater than 0")
	}
	if payment.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}
	if payment.NextPaymentDate.IsZero() {
		return fmt.Errorf("next payment date is required")
	}
	if payment.NextPaymentDate.Before(payment.DueDate) {
		return fmt.Errorf("next payment date must be after due date")
	}

	// Generate new ID if not provided
	if payment.ID == "" {
		payment.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	payment.CreatedAt = now
	payment.UpdatedAt = now

	return s.repo.Create(payment)
}

func (s *service) GetPayment(id string) (*Payment, error) {
	return s.repo.Get(id)
}

func (s *service) RecordPayment(id string, paidDate time.Time) error {
	payment, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	// Update payment with paid date
	payment.PaidDate = &paidDate
	payment.UpdatedAt = time.Now()

	return s.repo.Update(*payment)
}

func (s *service) UpdatePayment(payment Payment) error {
	// Verify payment exists
	existing, err := s.repo.Get(payment.ID)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	// Validate updates
	if payment.AmountDue <= 0 {
		return fmt.Errorf("amount due must be greater than 0")
	}
	if payment.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}
	if payment.NextPaymentDate.IsZero() {
		return fmt.Errorf("next payment date is required")
	}
	if payment.NextPaymentDate.Before(payment.DueDate) {
		return fmt.Errorf("next payment date must be after due date")
	}

	// Preserve original IDs and timestamps
	payment.TenantID = existing.TenantID
	payment.SpaceID = existing.SpaceID
	payment.CreatedAt = existing.CreatedAt
	payment.UpdatedAt = time.Now()

	return s.repo.Update(payment)
}

func (s *service) ListPaymentsByTenant(tenantID string) ([]Payment, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *service) ListPaymentsBySpace(spaceID string) ([]Payment, error) {
	return s.repo.ListBySpace(spaceID)
}

func (s *service) ListPaymentsByDate(date time.Time) ([]Payment, error) {
	return s.repo.ListByDate(date)
}

func (s *service) GetLatestPayment(spaceID string) (*Payment, error) {
	return s.repo.GetLatest(spaceID)
}
