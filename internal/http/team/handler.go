// Package team содержит обработчики и DTO для работы с командами.
package team

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/http/httperr"
	"github.com/dixitix/pr-reviewer-service/internal/service"
)

// Handler обрабатывает HTTP-запросы, связанные с командами.
type Handler struct {
	svc    service.TeamService
	logger *slog.Logger
}

// NewHandler создаёт новый обработчик команд.
func NewHandler(svc service.TeamService, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

// Add обрабатывает создание новой команды и её участников.
func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
		return
	}

	var req DTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeInvalidJSON, "invalid JSON body", h.logger)
		return
	}

	if req.TeamName == "" {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeValidation, "team_name is required", h.logger)
		return
	}

	for _, m := range req.Members {
		if m.UserID == "" || m.Username == "" {
			httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeValidation, "member.user_id and member.username are required", h.logger)
			return
		}
	}

	team, members := mapTeamDTOToDomain(req)

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Info("handleTeamAdd: creating team", slog.String("team_name", req.TeamName), slog.Int("members_count", len(req.Members)))
	}

	err := h.svc.CreateTeam(ctx, team.Name, members)
	if err != nil {
		if errors.Is(err, service.ErrTeamAlreadyExists) {
			httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeTeamExists, "team_name already exists", h.logger)
			return
		}

		if errors.Is(err, service.ErrNotFound) {
			httperr.WriteJSONError(w, http.StatusNotFound, httperr.ErrorCodeNotFound, "resource not found", h.logger)
			return
		}

		if h.logger != nil {
			h.logger.Error("handleTeamAdd: CreateTeam error", slog.Any("error", err))
		}

		httperr.WriteJSONError(w, http.StatusInternalServerError, httperr.ErrorCodeInternal, "internal server error", h.logger)
		return
	}

	resp := GetTeamResponse{
		Team: DTO{
			TeamName: req.TeamName,
			Members:  req.Members,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Error("handleTeamAdd: failed to write response", slog.Any("error", err))
		}
	}
}

// Get обрабатывает получение информации о команде.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httperr.WriteJSONError(w, http.StatusMethodNotAllowed, httperr.ErrorCodeMethodNotAllowed, "method not allowed", h.logger)
		return
	}

	teamNameParam := r.URL.Query().Get("team_name")
	if teamNameParam == "" {
		httperr.WriteJSONError(w, http.StatusBadRequest, httperr.ErrorCodeValidation, "team_name is required", h.logger)
		return
	}

	ctx := r.Context()
	team, members, err := h.svc.GetTeam(ctx, domain.TeamName(teamNameParam))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			httperr.WriteJSONError(w, http.StatusNotFound, httperr.ErrorCodeNotFound, "team not found", h.logger)
			return
		}

		if h.logger != nil {
			h.logger.Error("handleTeamGet: GetTeam error", slog.Any("error", err))
		}

		httperr.WriteJSONError(w, http.StatusInternalServerError, httperr.ErrorCodeInternal, "internal server error", h.logger)
		return
	}

	resp := mapTeamDomainToDTO(team, members)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Error("handleTeamGet: failed to write response", slog.Any("error", err))
		}
	}
}
