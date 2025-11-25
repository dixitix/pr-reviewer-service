// Package service содержит реализацию бизнес-логики сервиса назначения ревьюеров.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
)

// CreatePullRequest создаёт новый PR и назначает до двух ревьюверов из команды автора (исключая самого автора).
func (s *service) CreatePullRequest(
	ctx context.Context,
	id domain.PullRequestID,
	name string,
	authorID domain.UserID,
) (domain.PullRequest, error) {
	if _, err := s.pullRequestRepo.GetByID(ctx, id); err == nil {
		return domain.PullRequest{}, ErrPullRequestAlreadyExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return domain.PullRequest{}, fmt.Errorf("get pull request %s: %w", id, err)
	}

	// Получаем автора, чтобы узнать его команду.
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.PullRequest{}, ErrNotFound
		}

		return domain.PullRequest{}, fmt.Errorf("get author %s: %w", authorID, err)
	}

	// Проверяем, что команда автора существует.
	exists, err := s.teamRepo.TeamExists(ctx, author.TeamName)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("check team %s exists: %w", author.TeamName, err)
	}

	if !exists {
		return domain.PullRequest{}, ErrNotFound
	}

	// Берём активных участников команды автора, исключая самого автора.
	excludeAuthor := author.ID
	activeMembers, err := s.userRepo.ListActiveByTeam(ctx, author.TeamName, &excludeAuthor)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("list active users for team %s: %w", author.TeamName, err)
	}

	// Выбираем до двух ревьюверов из списка активных участников.
	reviewerIDs := s.pickReviewersForNewPR(activeMembers, 2)

	now := time.Now().UTC()

	pr := domain.PullRequest{
		ID:                id,
		Name:              name,
		AuthorID:          authorID,
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: reviewerIDs,
		CreatedAt:         &now,
		// MergedAt остаётся nil.
	}

	if err := s.pullRequestRepo.Create(ctx, pr); err != nil {
		return domain.PullRequest{}, fmt.Errorf("create pull request %s: %w", id, err)
	}

	return pr, nil
}

// MergePullRequest помечает PR как MERGED.
func (s *service) MergePullRequest(
	ctx context.Context,
	id domain.PullRequestID,
) (domain.PullRequest, error) {
	pr, err := s.pullRequestRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.PullRequest{}, ErrNotFound
		}

		return domain.PullRequest{}, fmt.Errorf("get pull request %s: %w", id, err)
	}

	if pr.Status == domain.PullRequestStatusMerged {
		// Идемпотентность: возвращаем текущее состояние без ошибок.
		return pr, nil
	}

	now := time.Now().UTC()
	pr.Status = domain.PullRequestStatusMerged
	pr.MergedAt = &now

	if err := s.pullRequestRepo.Update(ctx, pr); err != nil {
		return domain.PullRequest{}, fmt.Errorf("update pull request %s on merge: %w", id, err)
	}

	return pr, nil
}

// ReassignReviewer переназначает ревьювера на случайного активного участника из его команды.
func (s *service) ReassignReviewer(
	ctx context.Context,
	prID domain.PullRequestID,
	reviewerID domain.UserID,
) (domain.PullRequest, domain.UserID, error) {
	pr, err := s.pullRequestRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.PullRequest{}, "", ErrNotFound
		}

		return domain.PullRequest{}, "", fmt.Errorf("get pull request %s: %w", prID, err)
	}

	if pr.Status == domain.PullRequestStatusMerged {
		// После MERGED менять ревьюверов запрещено.
		return domain.PullRequest{}, "", ErrPullRequestMerged
	}

	// Проверяем, что пользователь действительно назначен ревьювером этого PR.
	reviewerIndex := -1
	for i, id := range pr.AssignedReviewers {
		if id == reviewerID {
			reviewerIndex = i
			break
		}
	}

	if reviewerIndex == -1 {
		return domain.PullRequest{}, "", ErrReviewerNotAssigned
	}

	reviewer, err := s.userRepo.GetByID(ctx, reviewerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.PullRequest{}, "", ErrNotFound
		}

		return domain.PullRequest{}, "", fmt.Errorf("get reviewer %s: %w", reviewerID, err)
	}

	// Берём активных участников команды заменяемого ревьювера.
	activeMembers, err := s.userRepo.ListActiveByTeam(ctx, reviewer.TeamName, nil)
	if err != nil {
		return domain.PullRequest{}, "", fmt.Errorf("list active users for team %s: %w", reviewer.TeamName, err)
	}

	candidates := s.buildReplacementCandidates(pr, reviewerID, activeMembers)
	if len(candidates) == 0 {
		return domain.PullRequest{}, "", ErrNoCandidate
	}

	s.rndMu.Lock()
	s.rnd.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	s.rndMu.Unlock()

	newReviewerID := candidates[0]

	pr.AssignedReviewers[reviewerIndex] = newReviewerID

	if err := s.pullRequestRepo.Update(ctx, pr); err != nil {
		return domain.PullRequest{}, "", fmt.Errorf("update pull request %s on reassign: %w", prID, err)
	}

	return pr, newReviewerID, nil
}
