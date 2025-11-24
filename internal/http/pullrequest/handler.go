// Package pullrequest содержит обработчики и DTO для работы с Pull Request'ами.
package pullrequest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/http/httperr"
	"github.com/dixitix/pr-reviewer-service/internal/service"
)

// Handler обрабатывает HTTP-запросы, связанные с Pull Request'ами.
type Handler struct {
	svc    service.PullRequestService
	logger *slog.Logger
}

// NewHandler создаёт обработчик для Pull Request-эндпоинтов.
func NewHandler(svc service.PullRequestService, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

// Create обрабатывает создание нового PR.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
		return
	}

	var req CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeInvalidJSON, "invalid JSON body", h.logger)
		return
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeValidation, "pull_request_id, pull_request_name and author_id are required", h.logger)
		return
	}

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(
			"handlePullRequestCreate",
			slog.String("pull_request_id", req.PullRequestID),
			slog.String("author_id", req.AuthorID),
		)
	}

	pr, err := h.svc.CreatePullRequest(ctx, domain.PullRequestID(req.PullRequestID), req.PullRequestName, domain.UserID(req.AuthorID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPullRequestAlreadyExists):
			httperr.WriteJSONError(w, http.StatusConflict, httperr.ErrorCodePRExists, "pull_request_id already exists", h.logger)
			return
		case errors.Is(err, service.ErrNotFound):
			httperr.WriteJSONError(w, http.StatusNotFound, httperr.ErrorCodeNotFound, "author or team not found", h.logger)
			return
		default:
			if h.logger != nil {
				h.logger.Error("handlePullRequestCreate: CreatePullRequest error", slog.Any("error", err))
			}
			httperr.WriteJSONError(w, http.StatusInternalServerError, httperr.ErrorCodeInternal, "internal server error", h.logger)
			return
		}
	}

	resp := Envelope{
		PullRequest: mapPullRequestDomainToDTO(pr),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Error("handlePullRequestCreate: failed to write response", slog.Any("error", err))
		}
	}
}

// Merge обрабатывает merge PR.
func (h *Handler) Merge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
		return
	}

	var req MergePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeInvalidJSON, "invalid JSON body", h.logger)
		return
	}

	if req.PullRequestID == "" {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeValidation, "pull_request_id is required", h.logger)
		return
	}

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info("handlePullRequestMerge", slog.String("pull_request_id", req.PullRequestID))
	}

	pr, err := h.svc.MergePullRequest(ctx, domain.PullRequestID(req.PullRequestID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httperr.WriteJSONError(w, http.StatusNotFound, httperr.ErrorCodeNotFound, "pull request not found", h.logger)
			return
		}

		if h.logger != nil {
			h.logger.Error("handlePullRequestMerge: MergePullRequest error", slog.Any("error", err))
		}

		httperr.WriteJSONError(w, http.StatusInternalServerError, httperr.ErrorCodeInternal, "internal server error", h.logger)
		return
	}

	resp := Envelope{
		PullRequest: mapPullRequestDomainToDTO(pr),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Error("handlePullRequestMerge: failed to write response", slog.Any("error", err))
		}
	}
}

// Reassign обрабатывает переназначение ревьювера.
func (h *Handler) Reassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
		return
	}

	var req ReassignPullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeInvalidJSON, "invalid JSON body", h.logger)
		return
	}

	if req.PullRequestID == "" || req.OldUserID == "" {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeValidation, "pull_request_id and old_user_id are required", h.logger)
		return
	}

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info(
			"handlePullRequestReassign",
			slog.String("pull_request_id", req.PullRequestID),
			slog.String("old_user_id", req.OldUserID),
		)
	}

	pr, newReviewerID, err := h.svc.ReassignReviewer(ctx, domain.PullRequestID(req.PullRequestID), domain.UserID(req.OldUserID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPullRequestMerged):
			httperr.WriteJSONError(w, http.StatusConflict, httperr.ErrorCodePRMerged, "pull request already merged", h.logger)
			return
		case errors.Is(err, service.ErrReviewerNotAssigned):
			httperr.WriteJSONError(w, http.StatusConflict, httperr.ErrorCodeNotAssigned, "old_user_id is not assigned as reviewer", h.logger)
			return
		case errors.Is(err, service.ErrNoCandidate):
			httperr.WriteJSONError(w, http.StatusConflict, httperr.ErrorCodeNoCandidate, "no candidate for reviewer reassignment", h.logger)
			return
		case errors.Is(err, service.ErrNotFound):
			httperr.WriteJSONError(w, http.StatusNotFound, httperr.ErrorCodeNotFound, "pull request or user not found", h.logger)
			return
		default:
			if h.logger != nil {
				h.logger.Error("handlePullRequestReassign: ReassignReviewer error", slog.Any("error", err))
			}
			httperr.WriteJSONError(w, http.StatusInternalServerError, httperr.ErrorCodeInternal, "internal server error", h.logger)
			return
		}
	}

	resp := ReassignResponse{
		PullRequest: mapPullRequestDomainToDTO(pr),
		ReplacedBy:  string(newReviewerID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Error("handlePullRequestReassign: failed to write response", slog.Any("error", err))
		}
	}
}
