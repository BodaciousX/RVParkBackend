// api/server.go
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
	s.Mux.Handle("POST /login", middleware.CORS(http.HandlerFunc(s.handleLogin)))

	// User routes - Admin only
	s.Mux.Handle("GET /users", middleware.CORS(authMiddleware.RequireAuth(authMiddleware.RequireAdmin(http.HandlerFunc(s.handleListUsers)))))
	s.Mux.Handle("POST /users", middleware.CORS(authMiddleware.RequireAuth(authMiddleware.RequireAdmin(http.HandlerFunc(s.handleCreateUser)))))
	s.Mux.Handle("GET /users/{id}", middleware.CORS(authMiddleware.RequireAuth(authMiddleware.RequireAdmin(http.HandlerFunc(s.handleGetUser)))))
	s.Mux.Handle("PUT /users/{id}", middleware.CORS(authMiddleware.RequireAuth(authMiddleware.RequireAdmin(http.HandlerFunc(s.handleUpdateUser)))))
	s.Mux.Handle("DELETE /users/{id}", middleware.CORS(authMiddleware.RequireAuth(authMiddleware.RequireAdmin(http.HandlerFunc(s.handleDeleteUser)))))

	// Space routes - Auth required
	s.Mux.Handle("GET /spaces", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleListSpaces))))
	s.Mux.Handle("GET /spaces/{id}", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleGetSpace))))
	s.Mux.Handle("PUT /spaces/{id}", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleUpdateSpace))))
	s.Mux.Handle("GET /spaces/vacant", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleGetVacantSpaces))))
	s.Mux.Handle("POST /spaces/{id}/reserve", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleReserveSpace))))
	s.Mux.Handle("POST /spaces/{id}/unreserve", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleUnreserveSpace))))
	s.Mux.Handle("POST /spaces/{id}/move-in", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleMoveIn))))
	s.Mux.Handle("POST /spaces/{id}/move-out", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleMoveOut))))

	// Tenant routes - Auth required
	s.Mux.Handle("GET /tenants", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleListTenants))))
	s.Mux.Handle("POST /tenants", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleCreateTenant))))
	s.Mux.Handle("GET /tenants/{id}", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleGetTenant))))
	s.Mux.Handle("PUT /tenants/{id}", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleUpdateTenant))))
	s.Mux.Handle("DELETE /tenants/{id}", middleware.CORS(authMiddleware.RequireAuth(authMiddleware.RequireAdmin(http.HandlerFunc(s.handleDeleteTenant)))))
	s.Mux.Handle("GET /tenants/{id}/payments", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleGetTenantPayments))))
	s.Mux.Handle("POST /tenants/{id}/payments", middleware.CORS(authMiddleware.RequireAuth(http.HandlerFunc(s.handleRecordPayment))))

	return s
}
