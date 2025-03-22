// api/tenant_handler.go contains the HTTP handlers for the tenant API.
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BodaciousX/RVParkBackend/tenant"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateTenantRequest struct {
	Name       string    `json:"name"`
	MoveInDate time.Time `json:"moveInDate"`
	SpaceID    string    `json:"spaceId"`
}

func (s *Server) handleListTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := s.tenantService.ListTenants()
	if err != nil {
		http.Error(w, "failed to fetch tenants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tenants); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCreateTenant(w http.ResponseWriter, r *http.Request) {
	var req CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newTenant := tenant.Tenant{
		ID:         uuid.New().String(),
		Name:       req.Name,
		MoveInDate: req.MoveInDate,
		SpaceID:    req.SpaceID,
	}

	if err := s.tenantService.CreateTenant(newTenant); err != nil {
		http.Error(w, "failed to create tenant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTenant)
}

func (s *Server) handleGetTenant(w http.ResponseWriter, r *http.Request) {
	// Extract id from URL path parameters
	vars := mux.Vars(r)
	id := vars["id"]

	tenant, err := s.tenantService.GetTenant(id)
	if err != nil {
		http.Error(w, "tenant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

func (s *Server) handleUpdateTenant(w http.ResponseWriter, r *http.Request) {
	// Extract id from URL path parameters
	vars := mux.Vars(r)
	id := vars["id"]

	var updateTenant tenant.Tenant
	if err := json.NewDecoder(r.Body).Decode(&updateTenant); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the ID in the path matches the tenant
	updateTenant.ID = id

	if err := s.tenantService.UpdateTenant(updateTenant); err != nil {
		http.Error(w, "failed to update tenant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateTenant)
}

func (s *Server) handleDeleteTenant(w http.ResponseWriter, r *http.Request) {
	// Extract id from URL path parameters
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.tenantService.DeleteTenant(id); err != nil {
		http.Error(w, "failed to delete tenant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
