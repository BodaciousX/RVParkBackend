// space/model.go contains the struct definitions for the space package.
package space

import (
	"time"
)

type Space struct {
	ID             string    `json:"id"`
	Section        string    `json:"section"` // e.g., "Mane Street"
	Status         string    `json:"status"`  // "Occupied (Paid)", "Occupied (Payment Due)", etc.
	TenantID       *string   `json:"tenantId,omitempty"`
	Reserved       bool      `json:"reserved"`
	PaymentType    string    `json:"paymentType,omitempty"` // "Monthly" or "Weekly"
	NextPayment    time.Time `json:"nextPayment,omitempty"`
	TenantNotified bool      `json:"tenantNotified"`
	PastDueAmount  float64   `json:"pastDueAmount"`
}
