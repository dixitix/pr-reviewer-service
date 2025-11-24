// Package user содержит обработчики и DTO для работы с пользователями.
package user

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/http/httperr"
	"github.com/dixitix/pr-reviewer-service/internal/http/pullrequest"
	"github.com/dixitix/pr-reviewer-service/internal/service"
)

// Handler обрабатывает HTTP-запросы, связанные с пользователями.
type Handler struct {
	svc    service.UserService
	logger *slog.Logger
}

// NewHandler создаёт обработчик пользователей.
func NewHandler(svc service.UserService, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

// SetIsActive обрабатывает изменение активности пользователя.
func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
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
		h.logger.Info("handleUserSetIsActive", slog.String("user_id", req.UserID), slog.Bool("is_active", req.IsActive))
	}

	user, err := h.svc.SetUserActive(ctx, domain.UserID(req.UserID), req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httperr.WriteJSONError(w, http.StatusNotFound, httperr.ErrorCodeNotFound, "user not found", h.logger)
			return
		}

		if h.logger != nil {
			h.logger.Error("handleUserSetIsActive: SetUserActive error", slog.Any("error", err))
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
			h.logger.Error("handleUserSetIsActive: failed to write response", slog.Any("error", err))
		}
	}
}

// GetReview обрабатывает получение списка PR'ов пользователя-ревьювера.
func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
		return
	}

	userIDParam := r.URL.Query().Get("user_id")
	if userIDParam == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info("handleUserGetReview", slog.String("user_id", userIDParam))
	}

	prs, err := h.svc.GetUserReviewPullRequests(ctx, domain.UserID(userIDParam))
	if err != nil && !errors.Is(err, service.ErrNotFound) {
		if h.logger != nil {
			h.logger.Error("handleUserGetReview: GetUserReviewPullRequests error", slog.Any("error", err))
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var prShort []pullrequest.Short
	if err == nil {
		prShort = pullrequest.MapPullRequestsToShort(prs)
	} else {
		prShort = []pullrequest.Short{}
	}

	resp := GetUserReviewResponse{
		UserID:       userIDParam,
		PullRequests: prShort,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Error("handleUserGetReview: failed to write response", slog.Any("error", err))
		}
	}
}
