// payment/model.go
package payment

import "time"

// PaymentMethod represents the method used for payment
type PaymentMethod string

// Constants for payment methods
const (
	PaymentMethodCredit PaymentMethod = "CREDIT"
	PaymentMethodCheck  PaymentMethod = "CHECK"
	PaymentMethodCash   PaymentMethod = "CASH"
)

// Payment represents a payment record for a tenant
type Payment struct {
	ID              string         `json:"id"`
	TenantID        string         `json:"tenantId"`
	AmountDue       float64        `json:"amountDue"`
	DueDate         time.Time      `json:"dueDate"`
	PaidDate        *time.Time     `json:"paidDate,omitempty"`
	NextPaymentDate time.Time      `json:"nextPaymentDate"`
	PaymentMethod   *PaymentMethod `json:"paymentMethod,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}
