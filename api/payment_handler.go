// api/payment_handler.go
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BodaciousX/RVParkBackend/payment"
	"github.com/google/uuid"
)

type CreatePaymentRequest struct {
	TenantID        string    `json:"tenantId"`
	AmountDue       float64   `json:"amountDue"`
	DueDate         time.Time `json:"dueDate"`
	NextPaymentDate time.Time `json:"nextPaymentDate"`
	PaidDate        time.Time `json:"paidDate"`
}

func (s *Server) handlePaymentList(w http.ResponseWriter, r *http.Request) {
	// Extract start and end dates from query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	tenantID := r.URL.Query().Get("tenant")

	if startStr == "" || endStr == "" {
		http.Error(w, "start and end dates are required", http.StatusBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "invalid start date format", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "invalid end date format", http.StatusBadRequest)
		return
	}

	var payments []payment.Payment
	if tenantID != "" {
		// If tenant ID is provided, get payments for specific tenant
		payments, err = s.paymentService.GetTenantPayments(tenantID)
	} else {
		// Otherwise get all payments in date range
		payments, err = s.paymentService.GetPaymentsByDateRange(start, end)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch payments: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payments); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.TenantID == "" {
		http.Error(w, "tenant ID is required", http.StatusBadRequest)
		return
	}
	if req.AmountDue <= 0 {
		http.Error(w, "amount due must be greater than 0", http.StatusBadRequest)
		return
	}
	if req.DueDate.IsZero() {
		http.Error(w, "due date is required", http.StatusBadRequest)
		return
	}
	if req.NextPaymentDate.IsZero() {
		http.Error(w, "next payment date is required", http.StatusBadRequest)
		return
	}

	newPayment := payment.Payment{
		ID:              uuid.New().String(),
		TenantID:        req.TenantID,
		AmountDue:       req.AmountDue,
		DueDate:         req.DueDate,
		PaidDate:        &req.PaidDate,
		NextPaymentDate: req.NextPaymentDate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.paymentService.CreatePayment(newPayment); err != nil {
		http.Error(w, fmt.Sprintf("failed to create payment: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newPayment); err != nil {
		// Log the error but don't return it to the client since we already sent the status code
		log.Printf("failed to encode response: %v", err)
	}
}

func (s *Server) handleGetPayment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/payments/")

	payment, err := s.paymentService.GetPayment(id)
	if err != nil {
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (s *Server) handleUpdatePayment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/payments/")

	var updatePayment payment.Payment
	if err := json.NewDecoder(r.Body).Decode(&updatePayment); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updatePayment.ID = id
	if err := s.paymentService.UpdatePayment(updatePayment); err != nil {
		http.Error(w, "failed to update payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatePayment)
}

func (s *Server) handleDeletePayment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/payments/")

	if err := s.paymentService.DeletePayment(id); err != nil {
		http.Error(w, "failed to delete payment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handlePaymentOperations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetPayment(w, r)
	case http.MethodPut:
		s.handleUpdatePayment(w, r)
	case http.MethodDelete:
		s.handleDeletePayment(w, r)
	default:
		http.NotFound(w, r)
	}
}
