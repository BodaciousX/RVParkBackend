// payment/interface.go
package payment

import "time"

type Service interface {
	// Core payment operations
	CreatePayment(payment Payment) error
	GetPayment(id string) (*Payment, error)
	UpdatePayment(payment Payment) error
	DeletePayment(id string) error

	// Query operations
	GetTenantPayments(tenantID string) ([]Payment, error)
	GetPaymentsByDateRange(start, end time.Time) ([]Payment, error)
	GetLatestPayment(tenantID string) (*Payment, error)
}

type Repository interface {
	// Core database operations
	Create(payment Payment) error
	Get(id string) (*Payment, error)
	Update(payment Payment) error
	Delete(id string) error

	// Query operations
	ListByTenant(tenantID string) ([]Payment, error)
	ListByDateRange(start, end time.Time) ([]Payment, error)
	GetLatestByTenant(tenantID string) (*Payment, error)
}
