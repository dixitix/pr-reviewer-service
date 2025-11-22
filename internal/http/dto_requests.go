// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

// SetUserActiveRequest описывает тело запроса /users/setIsActive.
type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// CreatePullRequestRequest описывает тело запроса /pullRequest/create.
type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

// MergePullRequestRequest описывает тело запроса /pullRequest/merge.
type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

// ReassignPullRequestRequest описывает тело запроса /pullRequest/reassign.
type ReassignPullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}
