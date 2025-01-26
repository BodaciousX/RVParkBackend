// space/interface.go contains the interface for the space package.
package space

type Service interface {
	ListSpaces() (map[string][]Space, error)
	GetSpace(id string) (*Space, error)
	GetVacantSpaces() ([]Space, error)
	ReserveSpace(spaceID string) error
	UnreserveSpace(spaceID string) error
	MoveIn(spaceID string, tenantID string) error
	MoveOut(spaceID string) error
}

type Repository interface {
	List() ([]Space, error)
	Get(id string) (*Space, error)
	Update(space Space) error
}
