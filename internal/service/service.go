// Package service содержит реализацию бизнес-логики сервиса назначения ревьюеров.
package service

import (
	"context"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
)

// TeamService описывает операции над командами и их участниками.
type TeamService interface {
	// CreateTeam создаёт новую команду и обновляет/создаёт её участников.
	CreateTeam(ctx context.Context, name domain.TeamName, members []domain.User) error

	// GetTeam возвращает команду и всех её участников.
	GetTeam(ctx context.Context, name domain.TeamName) (domain.Team, []domain.User, error)
}

// UserService описывает операции над пользователями.
type UserService interface {
	// SetUserActive меняет флаг активности пользователя и возвращает обновлённого пользователя.
	SetUserActive(ctx context.Context, userID domain.UserID, isActive bool) (domain.User, error)

	// GetUserReviewPullRequests возвращает список PR'ов, где пользователь выступает ревьювером.
	// Если пользователь не найден, возвращается ErrNotFound.
	GetUserReviewPullRequests(ctx context.Context, userID domain.UserID) ([]domain.PullRequest, error)
}

// PullRequestService описывает операции над Pull Request'ами.
type PullRequestService interface {
	// CreatePullRequest создаёт новый PR и назначает ревьюверов согласно правилам.
	CreatePullRequest(ctx context.Context, id domain.PullRequestID, name string, authorID domain.UserID) (domain.PullRequest, error)

	// MergePullRequest выполняет идемпотентный merge PR.
	MergePullRequest(ctx context.Context, id domain.PullRequestID) (domain.PullRequest, error)

	// ReassignReviewer переназначает ревьювера и возвращает обновлённый PR и user_id нового ревьювера.
	ReassignReviewer(ctx context.Context, prID domain.PullRequestID, oldReviewerID domain.UserID) (domain.PullRequest, domain.UserID, error)
}

// StatsService описывает операции получения статистики назначений.
type StatsService interface {
	// GetAssignmentsByUser возвращает количество назначений по каждому пользователю.
	GetAssignmentsByUser(ctx context.Context) (map[domain.UserID]int, error)

	// GetAssignmentsByPullRequest возвращает количество назначений по каждому Pull Request.
	GetAssignmentsByPullRequest(ctx context.Context) (map[domain.PullRequestID]int, error)
}

// UserAssignmentStat описывает количество назначений по пользователям.
type UserAssignmentStat struct {
	UserID      domain.UserID
	Assignments int
}

// PullRequestAssignmentStat описывает количество назначений по PR.
type PullRequestAssignmentStat struct {
	PullRequestID domain.PullRequestID
	Assignments   int
}

// Service агрегирует все доменные сервисы.
type Service interface {
	TeamService
	UserService
	PullRequestService
	StatsService
}
