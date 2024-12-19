// tenant/model.go contains the struct definitions for the tenant package.
package tenant

import "time"

type Tenant struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	MoveInDate time.Time `json:"moveInDate"`
	SpaceID    string    `json:"spaceId"`
}

type Payment struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenantId"`
	Amount      float64    `json:"amount"`
	DueDate     time.Time  `json:"dueDate"`
	PaidDate    *time.Time `json:"paidDate,omitempty"`
	PaymentType string     `json:"paymentType"` // "Monthly" or "Weekly"
	Status      string     `json:"status"`      // "Paid", "Due", "Overdue"
}
