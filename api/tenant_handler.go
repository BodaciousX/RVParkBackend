// api/tenant_handler.go contains the HTTP handlers for the tenant API.
package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/BodaciousX/RVParkBackend/tenant"
)

type CreateTenantRequest struct {
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	MoveInDate time.Time `json:"moveInDate"`
	SpaceID    string    `json:"spaceId"`
}

type RecordPaymentRequest struct {
	Amount      float64   `json:"amount"`
	DueDate     time.Time `json:"dueDate"`
	PaymentType string    `json:"paymentType"`
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
		Name:       req.Name,
		Phone:      req.Phone,
		MoveInDate: req.MoveInDate,
		SpaceID:    req.SpaceID,
	}

	if err := s.tenantService.CreateTenant(newTenant); err != nil {
		http.Error(w, "failed to create tenant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTenant)
}

func (s *Server) handleGetTenant(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tenants/")

	tenant, err := s.tenantService.GetTenant(id)
	if err != nil {
		http.Error(w, "tenant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

func (s *Server) handleUpdateTenant(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tenants/")

	var updateTenant tenant.Tenant
	if err := json.NewDecoder(r.Body).Decode(&updateTenant); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the ID in the path matches the tenant
	updateTenant.ID = id

	if err := s.tenantService.UpdateTenant(updateTenant); err != nil {
		http.Error(w, "failed to update tenant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateTenant)
}

func (s *Server) handleDeleteTenant(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tenants/")

	if err := s.tenantService.DeleteTenant(id); err != nil {
		http.Error(w, "failed to delete tenant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetTenantPayments(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/tenants/")
	tenantID := strings.TrimSuffix(path, "/payments")

	payments, err := s.tenantService.GetTenantPayments(tenantID)
	if err != nil {
		http.Error(w, "failed to get tenant payments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (s *Server) handleRecordPayment(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/tenants/")
	tenantID := strings.TrimSuffix(path, "/payments")

	var req RecordPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	payment := tenant.Payment{
		TenantID:    tenantID,
		Amount:      req.Amount,
		DueDate:     req.DueDate,
		PaymentType: req.PaymentType,
		Status:      "Paid", // Status will be set to Paid as this is a payment record
	}

	if err := s.tenantService.RecordPayment(payment); err != nil {
		http.Error(w, "failed to record payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}
