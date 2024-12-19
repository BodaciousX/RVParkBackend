// api/server.go contains the implementation of the API server.
package api

import (
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

	// Public routes
	s.Mux.Handle("POST /login", s.handleLogin())

	// User routes (protected)
	s.Mux.Handle("GET /users", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(user.RoleAdmin)(
			s.handleListUsers(),
		),
	))
	s.Mux.Handle("POST /users", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(user.RoleAdmin)(
			s.handleCreateUser(),
		),
	))
	s.Mux.Handle("GET /users/{id}", authMiddleware.RequireAuth(s.handleGetUser()))
	s.Mux.Handle("PUT /users/{id}", authMiddleware.RequireAuth(s.handleUpdateUser()))
	s.Mux.Handle("DELETE /users/{id}", authMiddleware.RequireAuth(
		authMiddleware.RequireRole(user.RoleAdmin)(
			s.handleDeleteUser(),
		),
	))

	// Space routes
	s.Mux.Handle("GET /spaces", authMiddleware.RequireAuth(s.handleListSpaces()))
	s.Mux.Handle("GET /spaces/{id}", authMiddleware.RequireAuth(s.handleGetSpace()))
	s.Mux.Handle("PUT /spaces/{id}", authMiddleware.RequireAuth(s.handleUpdateSpace()))
	s.Mux.Handle("GET /spaces/vacant", authMiddleware.RequireAuth(s.handleGetVacantSpaces()))
	s.Mux.Handle("POST /spaces/{id}/reserve", authMiddleware.RequireAuth(s.handleReserveSpace()))
	s.Mux.Handle("POST /spaces/{id}/unreserve", authMiddleware.RequireAuth(s.handleUnreserveSpace()))
	s.Mux.Handle("POST /spaces/{id}/move-in", authMiddleware.RequireAuth(s.handleMoveIn()))
	s.Mux.Handle("POST /spaces/{id}/move-out", authMiddleware.RequireAuth(s.handleMoveOut()))

	// Tenant routes
	s.Mux.Handle("GET /tenants", authMiddleware.RequireAuth(s.handleListTenants()))
	s.Mux.Handle("POST /tenants", authMiddleware.RequireAuth(s.handleCreateTenant()))
	s.Mux.Handle("GET /tenants/{id}", authMiddleware.RequireAuth(s.handleGetTenant()))
	s.Mux.Handle("PUT /tenants/{id}", authMiddleware.RequireAuth(s.handleUpdateTenant()))
	s.Mux.Handle("DELETE /tenants/{id}", authMiddleware.RequireAuth(s.handleDeleteTenant()))
	s.Mux.Handle("GET /tenants/{id}/payments", authMiddleware.RequireAuth(s.handleGetTenantPayments()))
	s.Mux.Handle("POST /tenants/{id}/payments", authMiddleware.RequireAuth(s.handleRecordPayment()))

	return s
}
