// Package stats содержит обработчики и DTO для выдачи статистики назначений.
package stats

// UserStat описывает количество назначений для конкретного пользователя.
type UserStat struct {
	UserID      string `json:"user_id"`
	Assignments int    `json:"assignments"`
}

// PullRequestStat описывает количество назначений для конкретного Pull Request.
type PullRequestStat struct {
	PullRequestID string `json:"pull_request_id"`
	Assignments   int    `json:"assignments"`
}

// UserStatsResponse описывает ответ на запрос о статистике по пользователю.
type UserStatsResponse struct {
	Stats []UserStat `json:"stats"`
}

// PullRequestStatsResponse описывает ответ на запрос о статистике по Pull Request'у.
type PullRequestStatsResponse struct {
	Stats []PullRequestStat `json:"stats"`
}
