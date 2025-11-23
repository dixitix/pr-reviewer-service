// Package service содержит реализацию бизнес-логики сервиса назначения ревьюеров.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
)

// SetUserActive меняет флаг активности пользователя и возвращает обновлённого пользователя.
func (s *service) SetUserActive(
	ctx context.Context,
	userID domain.UserID,
	isActive bool,
) (domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.User{}, ErrNotFound
		}

		return domain.User{}, fmt.Errorf("get user by id %s: %w", userID, err)
	}

	// Если флаг активности уже такой — операция идемпотентна.
	if user.IsActive == isActive {
		return user, nil
	}

	if err := s.userRepo.SetActive(ctx, userID, isActive); err != nil {
		return domain.User{}, fmt.Errorf("set user %s active=%t: %w", userID, isActive, err)
	}

	user.IsActive = isActive

	return user, nil
}

// GetUserReviewPullRequests возвращает список PR'ов, где пользователь выступает ревьювером.
func (s *service) GetUserReviewPullRequests(
	ctx context.Context,
	userID domain.UserID,
) ([]domain.PullRequest, error) {
	// Явно проверяем существование пользователя.
	// Это позволяет отличить "нет такого пользователя" от "нет PR'ов".
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get user by id %s: %w", userID, err)
	}

	prs, err := s.pullRequestRepo.ListByReviewer(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list pull requests for reviewer %s: %w", userID, err)
	}

	if prs == nil {
		return []domain.PullRequest{}, nil
	}

	return prs, nil
}
