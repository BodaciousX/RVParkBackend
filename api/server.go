// api/server.go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/BodaciousX/RVParkBackend/middleware"
	"github.com/BodaciousX/RVParkBackend/space"
	"github.com/BodaciousX/RVParkBackend/tenant"
	"github.com/BodaciousX/RVParkBackend/user"
)

type Server struct {
	Mux            *http.ServeMux
	userService    user.Service
	tenantService  tenant.Service
	spaceService   space.Service
	authMiddleware *middleware.AuthMiddleware
}

func NewServer(
	userService user.Service,
	tenantService tenant.Service,
	spaceService space.Service,
	authMiddleware *middleware.AuthMiddleware,
) *Server {
	s := &Server{
		Mux:            http.NewServeMux(),
		userService:    userService,
		tenantService:  tenantService,
		spaceService:   spaceService,
		authMiddleware: authMiddleware,
	}

	// Public routes with CORS
	s.Mux.Handle("/login", middleware.CORS(http.HandlerFunc(s.handleLogin)))
	s.Mux.Handle("/validate-token", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleValidateToken))))

	// Protected routes with auth and CORS
	// User routes - Admin only
	s.Mux.Handle("/users", middleware.CORS(authMiddleware.RequireAuth(authMiddleware.RequireAdmin(http.HandlerFunc(s.handleListUsers)))))
	s.Mux.Handle("/users/", middleware.CORS(authMiddleware.RequireAuth(authMiddleware.RequireAdmin(http.HandlerFunc(s.handleUserOperations)))))

	// Space routes
	s.Mux.Handle("/spaces", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleListSpaces))))
	s.Mux.Handle("/spaces/vacant", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleGetVacantSpaces))))
	s.Mux.Handle("/spaces/", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleSpaceOperations))))

	// Tenant routes
	s.Mux.Handle("/tenants", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleTenantList))))
	s.Mux.Handle("/tenants/", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleTenantOperations))))

	// Logout route
	s.Mux.Handle("/logout", middleware.CORS(
		authMiddleware.RequireAuth(
			http.HandlerFunc(s.handleLogout),
		),
	))

	return s
}

// Add new handler for validate-token
func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	// The user is already validated by the RequireAuth middleware
	// Just return the user from the context
	user := r.Context().Value(middleware.UserContextKey).(*user.User)

	// Return user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

// Handler for all user-related operations that need path parsing
func (s *Server) handleUserOperations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetUser(w, r)
	case http.MethodPut:
		s.handleUpdateUser(w, r)
	case http.MethodDelete:
		s.handleDeleteUser(w, r)
	default:
		http.NotFound(w, r)
	}
}

// Handler for all space-related operations that need path parsing
func (s *Server) handleSpaceOperations(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case r.Method == http.MethodGet:
		s.handleGetSpace(w, r)
	case r.Method == http.MethodPost && len(path) > 8:
		switch {
		case path[len(path)-8:] == "/reserve":
			s.handleReserveSpace(w, r)
		case path[len(path)-10:] == "/unreserve":
			s.handleUnreserveSpace(w, r)
		case path[len(path)-8:] == "/move-in":
			s.handleMoveIn(w, r)
		case path[len(path)-9:] == "/move-out":
			s.handleMoveOut(w, r)
		default:
			http.NotFound(w, r)
		}
	default:
		http.NotFound(w, r)
	}
}

// Handler for tenant list operations
func (s *Server) handleTenantList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleListTenants(w, r)
	case http.MethodPost:
		s.handleCreateTenant(w, r)
	default:
		http.NotFound(w, r)
	}
}

// Handler for all tenant-related operations that need path parsing
func (s *Server) handleTenantOperations(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case r.Method == http.MethodGet:
		if len(path) > 16 && path[len(path)-9:] == "/payments" {
			s.handleGetTenantPayments(w, r)
		} else {
			s.handleGetTenant(w, r)
		}
	case r.Method == http.MethodPut:
		s.handleUpdateTenant(w, r)
	case r.Method == http.MethodDelete:
		// Check if user is admin for delete operation
		if user, ok := r.Context().Value(middleware.UserContextKey).(*user.User); !ok || user.Role != "ADMIN" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		s.handleDeleteTenant(w, r)
	case r.Method == http.MethodPost && len(path) > 9 && path[len(path)-9:] == "/payments":
		s.handleRecordPayment(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*user.User)
	if err := s.userService.RevokeAllTokens(user.ID); err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
