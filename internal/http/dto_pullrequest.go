// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import "time"

// PullRequestDTO представляет полный Pull Request согласно OpenAPI-схеме.
type PullRequestDTO struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
}

// PullRequestShortDTO представляет укороченное представление PR для списков.
type PullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// PullRequestEnvelope оборачивает PullRequest в поле "pr".
type PullRequestEnvelope struct {
	PullRequest PullRequestDTO `json:"pr"`
}

// ReassignResponse описывает ответ на переназначение ревьювера.
type ReassignResponse struct {
	PullRequest PullRequestDTO `json:"pr"`
	ReplacedBy  string         `json:"replaced_by"`
}
