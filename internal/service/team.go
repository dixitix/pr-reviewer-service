// Package service содержит реализацию бизнес-логики сервиса назначения ревьюеров.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
)

// CreateTeam создаёт новую команду и при необходимости добавляет/обновляет её участников.
func (s *service) CreateTeam(
	ctx context.Context,
	name domain.TeamName,
	members []domain.User,
) error {
	team := domain.Team{
		Name: name,
	}

	if err := s.teamRepo.CreateTeam(ctx, team); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrTeamAlreadyExists
		}

		return fmt.Errorf("create team %s: %w", name, err)
	}

	if len(members) == 0 {
		return nil
	}

	if err := s.teamRepo.UpsertMembers(ctx, name, members); err != nil {
		return fmt.Errorf("upsert members for team %s: %w", name, err)
	}

	return nil
}

// GetTeam возвращает команду и всех её участников.
func (s *service) GetTeam(
	ctx context.Context,
	name domain.TeamName,
) (domain.Team, []domain.User, error) {
	team, members, err := s.teamRepo.GetTeamWithMembers(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Team{}, nil, ErrNotFound
		}

		return domain.Team{}, nil, fmt.Errorf("get team %s: %w", name, err)
	}

	return team, members, nil
}
