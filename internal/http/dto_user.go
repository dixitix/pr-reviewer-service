// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

// UserDTO представляет пользователя в HTTP-слое.
type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// SetUserActiveResponse описывает ответ на /users/setIsActive.
type SetUserActiveResponse struct {
	User UserDTO `json:"user"`
}

// GetUserReviewResponse описывает ответ на запрос списка PR'ов пользователя-ревьювера.
type GetUserReviewResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}
