// user/interface.go contains the interface for the user package.
package user

type Service interface {
	CreateUser(user User, password string) error
	GetUser(id string) (*User, error)
	UpdateUser(user User) error
	DeleteUser(id string) error
	Login(creds LoginCredentials) (*User, string, error)
	ValidateToken(token string) (*User, error)
	ChangePassword(userID string, oldPassword, newPassword string) error
}

type Repository interface {
	Create(user User) error
	Get(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user User) error
	Delete(id string) error
}
