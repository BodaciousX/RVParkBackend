// api/user_handler.go contains the HTTP handlers for the user API endpoints.
package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/BodaciousX/RVParkBackend/user"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  user.User `json:"user"`
	Token string    `json:"token"`
}

type CreateUserRequest struct {
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     user.Role `json:"role"`
}

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, token, err := s.userService.Login(user.LoginCredentials{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		resp := LoginResponse{
			User:  *user,
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *Server) handleListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The actual list users functionality would go here
		// For now, we'll return a placeholder response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}
}

func (s *Server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		newUser := user.User{
			Email:     req.Email,
			Username:  req.Username,
			Role:      req.Role,
			CreatedAt: time.Now(),
		}

		if err := s.userService.CreateUser(newUser, req.Password); err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}

func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/users/")

		user, err := s.userService.GetUser(id)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func (s *Server) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/users/")

		var updateUser user.User
		if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		updateUser.ID = id
		if err := s.userService.UpdateUser(updateUser); err != nil {
			http.Error(w, "failed to update user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updateUser)
	}
}

func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/users/")

		if err := s.userService.DeleteUser(id); err != nil {
			http.Error(w, "failed to delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
