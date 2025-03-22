// payment/p_interface.go
package payment

import "time"

type Service interface {
	CreatePayment(payment Payment) error
	GetPayment(id string) (*Payment, error)
	UpdatePayment(payment Payment) error
	DeletePayment(id string) error
	GetTenantPayments(tenantID string) ([]Payment, error)
	GetPaymentsByDateRange(start, end time.Time) ([]Payment, error)
	GetLatestPayment(tenantID string) (*Payment, error)
	RecordPayment(paymentID string, method PaymentMethod) error
}

type Repository interface {
	Create(payment Payment) error
	Get(id string) (*Payment, error)
	Update(payment Payment) error
	Delete(id string) error
	ListByTenant(tenantID string) ([]Payment, error)
	ListByDateRange(start, end time.Time) ([]Payment, error)
	ListByDateRangeAndTenant(start, end time.Time, tenantID string) ([]Payment, error)
	GetLatestByTenant(tenantID string) (*Payment, error)
}
