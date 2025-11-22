// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import "net/http"

// RegisterRoutes регистрирует HTTP-маршруты сервиса на переданном ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/team/add", h.handleTeamAdd)
	mux.HandleFunc("/team/get", h.handleTeamGet)
	mux.HandleFunc("/users/setIsActive", h.handleUserSetIsActive)
	mux.HandleFunc("/users/getReview", h.handleUserGetReview)
	mux.HandleFunc("/pullRequest/create", h.handlePullRequestCreate)
	mux.HandleFunc("/pullRequest/merge", h.handlePullRequestMerge)
	mux.HandleFunc("/pullRequest/reassign", h.handlePullRequestReassign)
}
