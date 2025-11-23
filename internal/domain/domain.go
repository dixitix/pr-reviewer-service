// Package domain содержит основные сущности сервиса назначения ревьюеров для Pull Request'ов.
package domain

import "time"

// UserID — тип идентификатора пользователя.
type UserID string

// TeamName — тип имени команды.
type TeamName string

// PullRequestID — тип идентификатора Pull Request'а.
type PullRequestID string

// PullRequestStatus описывает статус Pull Request'а в системе.
type PullRequestStatus string

const (
	// PullRequestStatusOpen — PR открыт и может изменяться.
	PullRequestStatusOpen PullRequestStatus = "OPEN"
	// PullRequestStatusMerged — PR смёржен, изменять его нельзя.
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

// User представляет пользователя, который может создавать PR и выступать ревьювером.
type User struct {
	ID       UserID
	Username string
	TeamName TeamName
	IsActive bool
}

// Team представляет команду разработчиков.
type Team struct {
	Name TeamName
}

// PullRequest представляет Pull Request и список назначенных ревьюверов.
type PullRequest struct {
	ID                PullRequestID
	Name              string
	AuthorID          UserID
	Status            PullRequestStatus
	AssignedReviewers []UserID
	CreatedAt         *time.Time
	MergedAt          *time.Time
}
