// payment/p_service.go
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
	if payment.AmountDue <= 0 {
		return fmt.Errorf("amount due must be greater than 0")
	}
	if payment.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
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

	// Preserve original IDs and timestamps
	payment.TenantID = existing.TenantID
	payment.CreatedAt = existing.CreatedAt
	payment.UpdatedAt = time.Now()

	return s.repo.Update(payment)
}

func (s *service) DeletePayment(id string) error {
	return s.repo.Delete(id)
}

func (s *service) GetTenantPayments(tenantID string) ([]Payment, error) {
	// Get a date range for the last 6 months
	end := time.Now()
	start := end.AddDate(0, -6, 0)
	return s.repo.ListByDateRangeAndTenant(start, end, tenantID)
}

func (s *service) GetPaymentsByDateRange(start, end time.Time) ([]Payment, error) {
	return s.repo.ListByDateRange(start, end)
}

func (s *service) GetLatestPayment(tenantID string) (*Payment, error) {
	return s.repo.GetLatestByTenant(tenantID)
}
