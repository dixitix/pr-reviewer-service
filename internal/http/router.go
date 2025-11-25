// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import "net/http"

// RegisterRoutes регистрирует HTTP-маршруты сервиса на переданном ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/team/add", h.teamHandler.Add)
	mux.HandleFunc("/team/get", h.teamHandler.Get)
	mux.HandleFunc("/users/setIsActive", h.userHandler.SetIsActive)
	mux.HandleFunc("/users/getReview", h.userHandler.GetReview)
	mux.HandleFunc("/pullRequest/create", h.pullRequestHandler.Create)
	mux.HandleFunc("/pullRequest/merge", h.pullRequestHandler.Merge)
	mux.HandleFunc("/pullRequest/reassign", h.pullRequestHandler.Reassign)
	mux.HandleFunc("/stats/byUser", h.statsHandler.AssignmentsByUser)
	mux.HandleFunc("/stats/byPullRequest", h.statsHandler.AssignmentsByPullRequest)
}
