// Package stats содержит обработчики и DTO для выдачи статистики назначений.
package stats

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dixitix/pr-reviewer-service/internal/http/httperr"
	"github.com/dixitix/pr-reviewer-service/internal/service"
)

// Handler обрабатывает HTTP-запросы, связанные со статистикой.
type Handler struct {
	svc    service.StatsService
	logger *slog.Logger
}

// NewHandler создаёт новый обработчик статистики.
func NewHandler(svc service.StatsService, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

// AssignmentsByUser отдаёт статистику назначений по пользователям.
func (h *Handler) AssignmentsByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
		return
	}

	stats, err := h.svc.GetAssignmentsByUser(r.Context())
	if err != nil {
		if h.logger != nil {
			h.logger.Error("handleStatsByUser: GetAssignmentsByUser error", slog.Any("error", err))
		}
		httperr.WriteJSONError(w, http.StatusInternalServerError, httperr.ErrorCodeInternal, "internal server error", h.logger)
		return
	}

	resp := UserStatsResponse{
		Stats: mapUserStatsToDTO(stats),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Error("handleStatsByUser: failed to write response", slog.Any("error", err))
		}
	}
}

// AssignmentsByPullRequest отдаёт статистику назначений по PR.
func (h *Handler) AssignmentsByPullRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
		return
	}

	stats, err := h.svc.GetAssignmentsByPullRequest(r.Context())
	if err != nil {
		if h.logger != nil {
			h.logger.Error("handleStatsByPullRequest: GetAssignmentsByPullRequest error", slog.Any("error", err))
		}
		httperr.WriteJSONError(w, http.StatusInternalServerError, httperr.ErrorCodeInternal, "internal server error", h.logger)
		return
	}

	resp := PullRequestStatsResponse{
		Stats: mapPullRequestStatsToDTO(stats),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Error("handleStatsByPullRequest: failed to write response", slog.Any("error", err))
		}
	}
}
