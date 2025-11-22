// Package repository содержит интерфейсы доступа к хранилищу данных
package repository

import (
	"context"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
)

// TeamRepository описывает операции с командами и их участниками.
type TeamRepository interface {
	// CreateTeam создаёт запись о команде.
	CreateTeam(ctx context.Context, team domain.Team) error

	// TeamExists возвращает true, если команда с таким именем существует.
	TeamExists(ctx context.Context, name domain.TeamName) (bool, error)

	// GetTeamWithMembers возвращает команду и всех её участников.
	GetTeamWithMembers(ctx context.Context, name domain.TeamName) (domain.Team, []domain.User, error)

	// UpsertMembers создаёт или обновляет пользователей команды по их ID.
	UpsertMembers(ctx context.Context, teamName domain.TeamName, members []domain.User) error
}

// UserRepository описывает операции с пользователями.
type UserRepository interface {
	// GetByID возвращает пользователя по ID.
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)

	// SetActive меняет флаг активности пользователя.
	SetActive(ctx context.Context, id domain.UserID, isActive bool) error

	// ListActiveByTeam возвращает активных пользователей команды.
	// Если excludeID != nil, пользователь с таким ID исключается из результата.
	ListActiveByTeam(ctx context.Context, teamName domain.TeamName, excludeID *domain.UserID) ([]domain.User, error)
}

// PullRequestRepository описывает операции с Pull Request'ами и их ревьюверами.
type PullRequestRepository interface {
	// Create создаёт новый PR вместе с назначенными ревьюверами.
	Create(ctx context.Context, pr domain.PullRequest) error

	// GetByID возвращает PR по его идентификатору.
	GetByID(ctx context.Context, id domain.PullRequestID) (domain.PullRequest, error)

	// Update обновляет состояние PR, включая список ревьюверов и статусы.
	Update(ctx context.Context, pr domain.PullRequest) error

	// ListByReviewer возвращает PR'ы, где пользователь является ревьювером.
	ListByReviewer(ctx context.Context, reviewerID domain.UserID) ([]domain.PullRequest, error)
}
