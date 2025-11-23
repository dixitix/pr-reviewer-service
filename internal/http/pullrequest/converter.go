// Package pullrequest содержит обработчики и DTO для работы с Pull Request'ами.
package pullrequest

import (
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
)

// mapPullRequestToShort конвертирует доменный PR в укороченный HTTP-DTO.
func mapPullRequestToShort(pr domain.PullRequest) Short {
	return Short{
		PullRequestID:   string(pr.ID),
		PullRequestName: pr.Name,
		AuthorID:        string(pr.AuthorID),
		Status:          string(pr.Status),
	}
}

// MapPullRequestsToShort конвертирует список доменных PR'ов в список укороченных DTO.
func MapPullRequestsToShort(prs []domain.PullRequest) []Short {
	result := make([]Short, len(prs))
	for i, pr := range prs {
		result[i] = mapPullRequestToShort(pr)
	}
	return result
}

// mapPullRequestDomainToDTO конвертирует доменный PR в полный HTTP-DTO.
func mapPullRequestDomainToDTO(pr domain.PullRequest) DTO {
	reviewers := make([]string, len(pr.AssignedReviewers))
	for i, id := range pr.AssignedReviewers {
		reviewers[i] = string(id)
	}

	var createdAt *time.Time
	if !pr.CreatedAt.IsZero() {
		t := pr.CreatedAt
		createdAt = t
	}

	var mergedAt *time.Time
	if pr.MergedAt != nil && !pr.MergedAt.IsZero() {
		t := *pr.MergedAt
		mergedAt = &t
	}

	return DTO{
		PullRequestID:     string(pr.ID),
		PullRequestName:   pr.Name,
		AuthorID:          string(pr.AuthorID),
		Status:            string(pr.Status),
		AssignedReviewers: reviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}
