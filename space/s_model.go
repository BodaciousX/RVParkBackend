// space/s_model.go
package space

type Space struct {
	ID       string  `json:"id"`
	Section  string  `json:"section"` // e.g., "Mane Street"
	Status   string  `json:"status"`  // "Occupied", "Vacant", "Reserved"
	TenantID *string `json:"tenantId,omitempty"`
	Reserved bool    `json:"reserved"`
}

// Constants for space status
const (
	StatusOccupied = "Occupied"
	StatusVacant   = "Vacant"
	StatusReserved = "Reserved"
)
