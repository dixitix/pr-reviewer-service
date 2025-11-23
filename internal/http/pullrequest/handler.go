// Package pullrequest содержит обработчики и DTO для работы с Pull Request'ами.
package pullrequest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/http/httperr"
	"github.com/dixitix/pr-reviewer-service/internal/service"
)

// Handler обрабатывает HTTP-запросы, связанные с Pull Request'ами.
type Handler struct {
	svc    service.PullRequestService
	logger *log.Logger
}

// NewHandler создаёт обработчик для Pull Request-эндпоинтов.
func NewHandler(svc service.PullRequestService, logger *log.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

// Create обрабатывает создание нового PR.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		http.Error(w, "pull_request_id, pull_request_name and author_id are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Printf("handlePullRequestCreate: pr_id=%s author_id=%s", req.PullRequestID, req.AuthorID)
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
				h.logger.Printf("handlePullRequestCreate: CreatePullRequest error: %v", err)
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
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
			h.logger.Printf("handlePullRequestCreate: failed to write response: %v", err)
		}
	}
}

// Merge обрабатывает merge PR.
func (h *Handler) Merge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MergePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" {
		http.Error(w, "pull_request_id is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Printf("handlePullRequestMerge: pr_id=%s", req.PullRequestID)
	}

	pr, err := h.svc.MergePullRequest(ctx, domain.PullRequestID(req.PullRequestID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httperr.WriteJSONError(w, http.StatusNotFound, httperr.ErrorCodeNotFound, "pull request not found", h.logger)
			return
		}

		if h.logger != nil {
			h.logger.Printf("handlePullRequestMerge: MergePullRequest error: %v", err)
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := Envelope{
		PullRequest: mapPullRequestDomainToDTO(pr),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Printf("handlePullRequestMerge: failed to write response: %v", err)
		}
	}
}

// Reassign обрабатывает переназначение ревьювера.
func (h *Handler) Reassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReassignPullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" || req.OldUserID == "" {
		http.Error(w, "pull_request_id and old_user_id are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Printf("handlePullRequestReassign: pr_id=%s old_user_id=%s", req.PullRequestID, req.OldUserID)
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
				h.logger.Printf("handlePullRequestReassign: ReassignReviewer error: %v", err)
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
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
			h.logger.Printf("handlePullRequestReassign: failed to write response: %v", err)
		}
	}
}
