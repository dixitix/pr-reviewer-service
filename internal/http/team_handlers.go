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

// handleTeamAdd обрабатывает создание новой команды и её участников.
func (h *Handler) handleTeamAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.TeamName == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}

	for _, m := range req.Members {
		if m.UserID == "" || m.Username == "" {
			http.Error(w, "member.user_id and member.username are required", http.StatusBadRequest)
			return
		}
	}

	team, members := mapTeamDTOToDomain(req)

	ctx := r.Context()

	if h.logger != nil {
		h.logger.Printf("handleTeamAdd: creating team %s with %d member(s)", req.TeamName, len(req.Members))
	}

	err := h.svc.CreateTeam(ctx, team.Name, members)
	if err != nil {
		if errors.Is(err, service.ErrTeamAlreadyExists) {
			writeJSONError(w, http.StatusBadRequest, errorCodeTeamExists, "team_name already exists")
			return
		}

		if errors.Is(err, service.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, errorCodeNotFound, "resource not found")
			return
		}

		if h.logger != nil {
			h.logger.Printf("handleTeamAdd: CreateTeam error: %v", err)
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := GetTeamResponse{
		Team: TeamDTO{
			TeamName: req.TeamName,
			Members:  req.Members,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Printf("handleTeamAdd: failed to write response: %v", err)
		}
	}
}

// handleTeamGet обрабатывает получение информации о команде.
func (h *Handler) handleTeamGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	teamNameParam := r.URL.Query().Get("team_name")
	if teamNameParam == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	team, members, err := h.svc.GetTeam(ctx, domain.TeamName(teamNameParam))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, errorCodeNotFound, "team not found")
			return
		}

		if h.logger != nil {
			h.logger.Printf("handleTeamGet: GetTeam error: %v", err)
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := mapTeamDomainToDTO(team, members)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if h.logger != nil {
			h.logger.Printf("handleTeamGet: failed to write response: %v", err)
		}
	}
}
