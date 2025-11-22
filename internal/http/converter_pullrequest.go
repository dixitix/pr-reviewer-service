// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import "github.com/dixitix/pr-reviewer-service/internal/domain"

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
