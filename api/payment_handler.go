// api/payment_handler.go
package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/BodaciousX/RVParkBackend/payment"
	"github.com/google/uuid"
)

type CreatePaymentRequest struct {
	TenantID string    `json:"tenantId"`
	Amount   float64   `json:"amount"`
	DueDate  time.Time `json:"dueDate"`
}

func (s *Server) handlePaymentList(w http.ResponseWriter, r *http.Request) {
	// Extract start and end dates from query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var payments []payment.Payment
	var err error

	if startStr != "" && endStr != "" {
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

		payments, err = s.paymentService.GetPaymentsByDateRange(start, end)
	} else {
		// If no date range provided, return an error as we want to enforce date range filtering
		http.Error(w, "start and end dates are required", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "failed to fetch payments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (s *Server) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newPayment := payment.Payment{
		ID:        uuid.New().String(),
		TenantID:  req.TenantID,
		AmountDue: req.Amount,
		DueDate:   req.DueDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.paymentService.CreatePayment(newPayment); err != nil {
		http.Error(w, "failed to create payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPayment)
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
