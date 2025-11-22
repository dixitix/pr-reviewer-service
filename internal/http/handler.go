// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import (
	"log"

	"github.com/dixitix/pr-reviewer-service/internal/service"
)

// Handler обрабатывает HTTP-запросы согласно OpenAPI-спецификации.
type Handler struct {
	svc    service.Service
	logger *log.Logger
}

// NewHandler создаёт новый HTTP-обработчик.
func NewHandler(svc service.Service, logger *log.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}
