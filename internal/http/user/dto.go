// Package user содержит обработчики и DTO для работы с пользователями.
package user

import "github.com/dixitix/pr-reviewer-service/internal/http/pullrequest"

// DTO представляет пользователя в HTTP-слое.
type DTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// SetUserActiveResponse описывает ответ на /users/setIsActive.
type SetUserActiveResponse struct {
	User DTO `json:"user"`
}

// GetUserReviewResponse описывает ответ на запрос списка PR'ов пользователя-ревьювера.
type GetUserReviewResponse struct {
	UserID       string              `json:"user_id"`
	PullRequests []pullrequest.Short `json:"pull_requests"`
}
