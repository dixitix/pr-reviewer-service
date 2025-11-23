// Package pullrequest содержит обработчики и DTO для работы с Pull Request'ами.
package pullrequest

import "time"

// DTO представляет полный Pull Request.
type DTO struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
}

// Short представляет укороченное представление PR для списков.
type Short struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// Envelope оборачивает PullRequest в поле "pr".
type Envelope struct {
	PullRequest DTO `json:"pr"`
}

// ReassignResponse описывает ответ на переназначение ревьювера.
type ReassignResponse struct {
	PullRequest DTO    `json:"pr"`
	ReplacedBy  string `json:"replaced_by"`
}
