// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import (
	"log/slog"

	"github.com/dixitix/pr-reviewer-service/internal/http/pullrequest"
	"github.com/dixitix/pr-reviewer-service/internal/http/stats"
	"github.com/dixitix/pr-reviewer-service/internal/http/team"
	"github.com/dixitix/pr-reviewer-service/internal/http/user"
	"github.com/dixitix/pr-reviewer-service/internal/service"
)

// Handler агрегирует обработчики HTTP-запросов.
type Handler struct {
	teamHandler        *team.Handler
	userHandler        *user.Handler
	pullRequestHandler *pullrequest.Handler
	statsHandler       *stats.Handler
}

// NewHandler создаёт новый HTTP-обработчик.
func NewHandler(svc service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		teamHandler:        team.NewHandler(svc, logger),
		userHandler:        user.NewHandler(svc, logger),
		pullRequestHandler: pullrequest.NewHandler(svc, logger),
		statsHandler:       stats.NewHandler(svc, logger),
	}
}
