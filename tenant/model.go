// tenant/model.go contains the struct definitions for the tenant package.
package tenant

import "time"

const (
	PaymentTypeMonthly = "Monthly"
	PaymentTypeWeekly  = "Weekly"
	PaymentTypeDaily   = "Daily"
)

const (
	PaymentStatusPaid    = "Paid"
	PaymentStatusDue     = "Due"
	PaymentStatusOverdue = "Overdue"
)

type Tenant struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	MoveInDate time.Time `json:"moveInDate"`
	SpaceID    string    `json:"spaceId"`
}

type Payment struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenantId"`
	Amount      float64    `json:"amount"`
	DueDate     time.Time  `json:"dueDate"`
	PaidDate    *time.Time `json:"paidDate,omitempty"`
	PaymentType string     `json:"paymentType"` // PaymentTypeMonthly, PaymentTypeWeekly, PaymentTypeDaily
	Status      string     `json:"status"`      // PaymentStatusPaid, PaymentStatusDue, PaymentStatusOverdue
}
