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

// handlePullRequestCreate обрабатывает создание нового PR.
func (h *Handler) handlePullRequestCreate(w http.ResponseWriter, r *http.Request) {
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
		h.logger.Printf(
			"handlePullRequestCreate: pr_id=%s author_id=%s",
			req.PullRequestID,
			req.AuthorID,
		)
	}

	pr, err := h.svc.CreatePullRequest(
		ctx,
		domain.PullRequestID(req.PullRequestID),
		req.PullRequestName,
		domain.UserID(req.AuthorID),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPullRequestAlreadyExists):
			writeJSONError(w, http.StatusConflict, errorCodePRExists, "pull_request_id already exists")
			return

		case errors.Is(err, service.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, errorCodeNotFound, "author or team not found")
			return

		default:
			if h.logger != nil {
				h.logger.Printf("handlePullRequestCreate: CreatePullRequest error: %v", err)
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp := PullRequestEnvelope{
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

// handlePullRequestMerge обрабатывает merge PR.
func (h *Handler) handlePullRequestMerge(w http.ResponseWriter, r *http.Request) {
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
			writeJSONError(w, http.StatusNotFound, errorCodeNotFound, "pull request not found")
			return
		}

		if h.logger != nil {
			h.logger.Printf("handlePullRequestMerge: MergePullRequest error: %v", err)
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := PullRequestEnvelope{
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

// handlePullRequestReassign обрабатывает переназначение ревьювера.
func (h *Handler) handlePullRequestReassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO handling route
}
