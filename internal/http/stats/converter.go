// Package stats содержит обработчики и DTO для выдачи статистики назначений.
package stats

import (
	"sort"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
)

// mapUserStatsToDTO конвертирует map статистики по пользователям в отсортированный список DTO.
func mapUserStatsToDTO(stats map[domain.UserID]int) []UserStat {
	if len(stats) == 0 {
		return []UserStat{}
	}

	ids := make([]string, 0, len(stats))
	for id := range stats {
		ids = append(ids, string(id))
	}

	sort.Strings(ids)

	result := make([]UserStat, 0, len(ids))
	for _, id := range ids {
		result = append(result, UserStat{
			UserID:      id,
			Assignments: stats[domain.UserID(id)],
		})
	}

	return result
}

// mapPullRequestStatsToDTO конвертирует map статистики по PR в отсортированный список DTO.
func mapPullRequestStatsToDTO(stats map[domain.PullRequestID]int) []PullRequestStat {
	if len(stats) == 0 {
		return []PullRequestStat{}
	}

	ids := make([]string, 0, len(stats))
	for id := range stats {
		ids = append(ids, string(id))
	}

	sort.Strings(ids)

	result := make([]PullRequestStat, 0, len(ids))
	for _, id := range ids {
		result = append(result, PullRequestStat{
			PullRequestID: id,
			Assignments:   stats[domain.PullRequestID(id)],
		})
	}

	return result
}
