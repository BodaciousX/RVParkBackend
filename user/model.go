// user/model.go contains the struct definitions for the user package.
package user

import "time"

type Role string

const (
	RoleAdmin   Role = "ADMIN"
	RoleManager Role = "MANAGER"
	RoleStaff   Role = "STAFF"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	LastLogin    time.Time `json:"lastLogin"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
