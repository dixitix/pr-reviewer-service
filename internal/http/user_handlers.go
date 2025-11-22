// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/service"
)

// handleUserSetIsActive обрабатывает изменение активности пользователя.
func (h *Handler) handleUserSetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Printf("handleUserSetIsActive: user_id=%s is_active=%v", req.UserID, req.IsActive)
	}

	user, err := h.svc.SetUserActive(ctx, domain.UserID(req.UserID), req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, errorCodeNotFound, "user not found")
			return
		}

		if h.logger != nil {
			h.logger.Printf("handleUserSetIsActive: SetUserActive error: %v", err)
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := SetUserActiveResponse{
		User: mapUserDomainToDTO(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Printf("handleUserSetIsActive: failed to write response: %v", err)
		}
	}
}

// handleUserGetReview обрабатывает получение списка PR'ов,
func (h *Handler) handleUserGetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO handling route
}
