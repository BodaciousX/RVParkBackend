// space_handler.go contains the HTTP handlers for the space API.
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/BodaciousX/RVParkBackend/space"
)

func (s *Server) handleListSpaces(w http.ResponseWriter, r *http.Request) {
	spaces, err := s.spaceService.ListSpaces()
	if err != nil {
		http.Error(w, "failed to list spaces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spaces)
}

func (s *Server) handleUpdateSpace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/spaces/")

	var updateSpace space.Space
	if err := json.NewDecoder(r.Body).Decode(&updateSpace); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the ID in the path matches the space
	updateSpace.ID = id

	// Get current space first
	currentSpace, err := s.spaceService.GetSpace(id)
	if err != nil {
		http.Error(w, "space not found", http.StatusNotFound)
		return
	}

	// Preserve the section from the current space
	updateSpace.Section = currentSpace.Section

	// Update the space
	if err := s.spaceService.UpdateSpace(updateSpace); err != nil {
		http.Error(w, "failed to update space", http.StatusInternalServerError)
		return
	}

	// Get the updated space to return
	updatedSpace, err := s.spaceService.GetSpace(id)
	if err != nil {
		http.Error(w, "failed to get updated space", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSpace)
}

func (s *Server) handleGetSpace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/spaces/")

	space, err := s.spaceService.GetSpace(id)
	if err != nil {
		http.Error(w, "space not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(space)
}

func (s *Server) handleGetVacantSpaces(w http.ResponseWriter, r *http.Request) {
	spaces, err := s.spaceService.GetVacantSpaces()
	if err != nil {
		http.Error(w, "failed to get vacant spaces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spaces)
}

func (s *Server) handleReserveSpace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/spaces/")
	id = strings.TrimSuffix(id, "/reserve")

	if err := s.spaceService.ReserveSpace(id); err != nil {
		http.Error(w, "failed to reserve space", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUnreserveSpace(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/spaces/")
	id = strings.TrimSuffix(id, "/unreserve")

	if err := s.spaceService.UnreserveSpace(id); err != nil {
		http.Error(w, "failed to unreserve space", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type MoveInRequest struct {
	TenantID string `json:"tenantId"`
}

func (s *Server) handleMoveIn(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/spaces/")
	id = strings.TrimSuffix(id, "/move-in")

	var req MoveInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.spaceService.MoveIn(id, req.TenantID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated space to return
	updatedSpace, err := s.spaceService.GetSpace(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSpace)
}

func (s *Server) handleMoveOut(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/spaces/")
	id = strings.TrimSuffix(id, "/move-out")

	if err := s.spaceService.MoveOut(id); err != nil {
		http.Error(w, "failed to move out tenant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
