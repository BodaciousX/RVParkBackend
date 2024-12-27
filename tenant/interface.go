// tenant/interface.go
package tenant

type Service interface {
	ListTenants() ([]Tenant, error)
	CreateTenant(tenant Tenant) error
	GetTenant(id string) (*Tenant, error)
	UpdateTenant(tenant Tenant) error
	DeleteTenant(id string) error
	GetTenantPayments(tenantID string) ([]Payment, error)
	RecordPayment(payment Payment) error
	GetPaymentStatus(tenantID string) (string, float64, error)
}

type Repository interface {
	List() ([]Tenant, error)
	Create(tenant Tenant) error
	Get(id string) (*Tenant, error)
	Update(tenant Tenant) error
	Delete(id string) error
	ListPayments(tenantID string) ([]Payment, error)
	CreatePayment(payment Payment) error
	UpdatePayment(payment Payment) error
}
