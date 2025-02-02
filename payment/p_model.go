// payment/p_model.go
package payment

import "time"

// Payment represents a payment record for a tenant
type Payment struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenantId"`
	AmountDue       float64    `json:"amountDue"`
	DueDate         time.Time  `json:"dueDate"`
	PaidDate        *time.Time `json:"paidDate,omitempty"`
	NextPaymentDate time.Time  `json:"nextPaymentDate"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
