// tenant/t_model.go
package tenant

import "time"

// Tenant represents a resident of the RV park
type Tenant struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	MoveInDate time.Time `json:"moveInDate"`
	SpaceID    string    `json:"spaceId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
