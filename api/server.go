// api/server.go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/BodaciousX/RVParkBackend/middleware"
	"github.com/BodaciousX/RVParkBackend/payment"
	"github.com/BodaciousX/RVParkBackend/space"
	"github.com/BodaciousX/RVParkBackend/tenant"
	"github.com/BodaciousX/RVParkBackend/user"
	"github.com/gorilla/mux"
)

type Server struct {
	Router         *mux.Router
	userService    user.Service
	tenantService  tenant.Service
	spaceService   space.Service
	paymentService payment.Service
	authMiddleware *middleware.AuthMiddleware
}

func NewServer(
	userService user.Service,
	tenantService tenant.Service,
	spaceService space.Service,
	paymentService payment.Service,
	authMiddleware *middleware.AuthMiddleware,
) *Server {
	s := &Server{
		Router:         mux.NewRouter(),
		userService:    userService,
		tenantService:  tenantService,
		spaceService:   spaceService,
		paymentService: paymentService,
		authMiddleware: authMiddleware,
	}

	// Apply CORS middleware to all routes
	s.Router.Use(middleware.CORSMiddleware)

	// Public routes
	s.Router.HandleFunc("/login", s.handleLogin).Methods("POST")
	s.Router.Handle("/validate-token", authMiddleware.RequireAuth(http.HandlerFunc(s.handleValidateToken))).Methods("GET")

	// User routes
	users := s.Router.PathPrefix("/users").Subrouter()
	users.Use(authMiddleware.RequireAuth)
	users.HandleFunc("", s.handleListUsers).Methods("GET")
	users.HandleFunc("", s.handleCreateUser).Methods("POST")
	users.HandleFunc("/{id}", s.handleGetUser).Methods("GET")
	users.HandleFunc("/{id}", s.handleUpdateUser).Methods("PUT")
	users.HandleFunc("/{id}", s.handleDeleteUser).Methods("DELETE")

	// Space routes
	spaces := s.Router.PathPrefix("/spaces").Subrouter()
	spaces.Use(authMiddleware.RequireAuth)
	spaces.HandleFunc("", s.handleListSpaces).Methods("GET")
	spaces.HandleFunc("/vacant", s.handleGetVacantSpaces).Methods("GET")
	spaces.HandleFunc("/{id}", s.handleGetSpace).Methods("GET")
	spaces.HandleFunc("/{id}", s.handleUpdateSpace).Methods("PUT")
	spaces.HandleFunc("/{id}/reserve", s.handleReserveSpace).Methods("POST")
	spaces.HandleFunc("/{id}/unreserve", s.handleUnreserveSpace).Methods("POST")
	spaces.HandleFunc("/{id}/move-in", s.handleMoveIn).Methods("POST")
	spaces.HandleFunc("/{id}/move-out", s.handleMoveOut).Methods("POST")

	// Tenant routes
	tenants := s.Router.PathPrefix("/tenants").Subrouter()
	tenants.Use(authMiddleware.RequireAuth)
	tenants.HandleFunc("", s.handleListTenants).Methods("GET")
	tenants.HandleFunc("", s.handleCreateTenant).Methods("POST")
	tenants.HandleFunc("/{id}", s.handleGetTenant).Methods("GET")
	tenants.HandleFunc("/{id}", s.handleUpdateTenant).Methods("PUT")
	tenants.HandleFunc("/{id}", s.handleDeleteTenant).Methods("DELETE")

	// Payment routes
	payments := s.Router.PathPrefix("/payments").Subrouter()
	payments.Use(authMiddleware.RequireAuth)
	payments.HandleFunc("", s.handlePaymentList).Methods("GET")
	payments.HandleFunc("", s.handleCreatePayment).Methods("POST")
	payments.HandleFunc("/{id}", s.handleGetPayment).Methods("GET")
	payments.HandleFunc("/{id}", s.handleUpdatePayment).Methods("PUT")
	payments.HandleFunc("/{id}", s.handleDeletePayment).Methods("DELETE")
	payments.HandleFunc("/{id}/record", s.handleRecordPayment).Methods("POST")

	// Logout route
	s.Router.Handle("/logout", authMiddleware.RequireAuth(http.HandlerFunc(s.handleLogout))).Methods("POST")

	return s
}

func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*user.User)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*user.User)
	if err := s.userService.RevokeAllTokens(user.ID); err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ServeHTTP makes the server struct satisfy http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
