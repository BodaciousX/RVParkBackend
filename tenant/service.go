// tenant/service.go contains the implementation of the Service interface.
package tenant

import (
	"time"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTenant(tenant Tenant) error {
	return s.repo.Create(tenant)
}

func (s *service) GetTenant(id string) (*Tenant, error) {
	return s.repo.Get(id)
}

func (s *service) UpdateTenant(tenant Tenant) error {
	return s.repo.Update(tenant)
}

func (s *service) DeleteTenant(id string) error {
	return s.repo.Delete(id)
}

func (s *service) ListTenants() ([]Tenant, error) {
	return s.repo.List()
}

func (s *service) GetTenantPayments(tenantID string) ([]Payment, error) {
	return s.repo.ListPayments(tenantID)
}

func (s *service) RecordPayment(payment Payment) error {
	now := time.Now()
	payment.PaidDate = &now
	payment.Status = "Paid"
	return s.repo.CreatePayment(payment)
}

func (s *service) GetPaymentStatus(tenantID string) (string, float64, error) {
	payments, err := s.repo.ListPayments(tenantID)
	if err != nil {
		return "", 0, err
	}

	var totalPastDue float64
	for _, payment := range payments {
		if payment.Status == "Overdue" {
			totalPastDue += payment.Amount
		}
	}

	if totalPastDue > 0 {
		return "Overdue", totalPastDue, nil
	}
	return "Paid", 0, nil
}
