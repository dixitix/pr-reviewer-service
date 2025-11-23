// Package service описывает интерфейсы и доменные ошибки бизнес-логики сервиса назначения ревьюеров.
package service

import "errors"

// Базовые доменные ошибки сервиса.
var (
	ErrTeamAlreadyExists        = errors.New("team already exists")
	ErrNotFound                 = errors.New("resource not found")
	ErrPullRequestAlreadyExists = errors.New("pull request already exists")
	ErrPullRequestMerged        = errors.New("pull request already merged")
	ErrReviewerNotAssigned      = errors.New("reviewer is not assigned to pull request")
	ErrNoCandidate              = errors.New("no candidate for reviewer reassignment")
)
