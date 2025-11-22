// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import (
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
)

// mapPullRequestToShortDTO конвертирует доменный PR в укороченный HTTP-DTO.
func mapPullRequestToShortDTO(pr domain.PullRequest) PullRequestShortDTO {
	return PullRequestShortDTO{
		PullRequestID:   string(pr.ID),
		PullRequestName: pr.Name,
		AuthorID:        string(pr.AuthorID),
		Status:          string(pr.Status),
	}
}

// mapPullRequestsToShortDTOs конвертирует список доменных PR'ов в список укороченных DTO.
func mapPullRequestsToShortDTOs(prs []domain.PullRequest) []PullRequestShortDTO {
	result := make([]PullRequestShortDTO, len(prs))
	for i, pr := range prs {
		result[i] = mapPullRequestToShortDTO(pr)
	}
	return result
}

// mapPullRequestDomainToDTO конвертирует доменный PR в полный HTTP-DTO.
func mapPullRequestDomainToDTO(pr domain.PullRequest) PullRequestDTO {
	reviewers := make([]string, len(pr.AssignedReviewers))
	for i, id := range pr.AssignedReviewers {
		reviewers[i] = string(id)
	}

	var createdAt *time.Time
	if !pr.CreatedAt.IsZero() {
		t := pr.CreatedAt
		createdAt = &t
	}

	var mergedAt *time.Time
	if pr.MergedAt != nil && !pr.MergedAt.IsZero() {
		t := *pr.MergedAt
		mergedAt = &t
	}

	return PullRequestDTO{
		PullRequestID:     string(pr.ID),
		PullRequestName:   pr.Name,
		AuthorID:          string(pr.AuthorID),
		Status:            string(pr.Status),
		AssignedReviewers: reviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}
