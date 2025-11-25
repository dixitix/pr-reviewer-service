// Package service содержит реализацию бизнес-логики сервиса назначения ревьюеров.
package service

import (
	"context"
	"fmt"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
)

// GetAssignmentsByUser возвращает количество назначений по каждому пользователю.
func (s *service) GetAssignmentsByUser(
	ctx context.Context,
) (map[domain.UserID]int, error) {
	stats, err := s.pullRequestRepo.CountAssignmentsByReviewer(ctx)
	if err != nil {
		return nil, fmt.Errorf("count assignments by reviewer: %w", err)
	}

	if stats == nil {
		return map[domain.UserID]int{}, nil
	}

	return stats, nil
}

// GetAssignmentsByPullRequest возвращает количество назначений по каждому Pull Request.
func (s *service) GetAssignmentsByPullRequest(
	ctx context.Context,
) (map[domain.PullRequestID]int, error) {
	stats, err := s.pullRequestRepo.CountAssignmentsByPullRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("count assignments by pull request: %w", err)
	}

	if stats == nil {
		return map[domain.PullRequestID]int{}, nil
	}

	return stats, nil
}
