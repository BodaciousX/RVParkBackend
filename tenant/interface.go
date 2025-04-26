// tenant/interface.go
package tenant

// TenantService interface now indicates it depends on SpaceService
type Service interface {
	// Core tenant management
	ListTenants() ([]Tenant, error)
	CreateTenant(tenant Tenant) error
	GetTenant(id string) (*Tenant, error)
	UpdateTenant(tenant Tenant) error
	DeleteTenant(id string) error

	// Utility methods
	GetTenantBySpace(spaceID string) (*Tenant, error)
}

type Repository interface {
	// Core tenant management
	List() ([]Tenant, error)
	Create(tenant Tenant) error
	Get(id string) (*Tenant, error)
	Update(tenant Tenant) error
	Delete(id string) error

	// Additional queries
	GetBySpace(spaceID string) (*Tenant, error)
}
