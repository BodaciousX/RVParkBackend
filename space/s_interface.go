// space/s_interface.go
package space

// No import for tenant package needed

type Service interface {
	ListSpaces() (map[string][]Space, error)
	GetSpace(id string) (*Space, error)
	GetVacantSpaces() ([]Space, error)
	ReserveSpace(spaceID string) error
	UnreserveSpace(spaceID string) error
	MoveIn(spaceID string, tenantID string) error
	MoveOut(spaceID string) error
	UpdateSpace(space Space) error
}

type Repository interface {
	List() ([]Space, error)
	Get(id string) (*Space, error)
	Update(space Space) error
}
