// Package postgres содержит реализацию репозиториев поверх PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
)

// TeamRepository реализует repository.TeamRepository поверх *sql.DB.
type TeamRepository struct {
	db *sql.DB
}

// NewTeamRepository создаёт новый экземпляр TeamRepository.
func NewTeamRepository(db *sql.DB) repository.TeamRepository {
	return &TeamRepository{db: db}
}

// CreateTeam создаёт запись о команде.
func (r *TeamRepository) CreateTeam(ctx context.Context, team domain.Team) error {
	exists, err := r.TeamExists(ctx, team.Name)
	if err != nil {
		return fmt.Errorf("check team %s exists: %w", team.Name, err)
	}

	if exists {
		return repository.ErrAlreadyExists
	}

	const query = `
		INSERT INTO teams (name)
		VALUES ($1)
	`

	if _, err := r.db.ExecContext(ctx, query, string(team.Name)); err != nil {
		return fmt.Errorf("insert team %s: %w", team.Name, err)
	}

	return nil
}

// TeamExists возвращает true, если команда с таким именем существует.
func (r *TeamRepository) TeamExists(ctx context.Context, name domain.TeamName) (bool, error) {
	const query = `
		SELECT 1
		FROM teams
		WHERE name = $1
	`

	var dummy int
	err := r.db.QueryRowContext(ctx, query, string(name)).Scan(&dummy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("query team %s exists: %w", name, err)
	}

	return true, nil
}

// GetTeamWithMembers возвращает команду и всех её участников.
func (r *TeamRepository) GetTeamWithMembers(
	ctx context.Context,
	name domain.TeamName,
) (domain.Team, []domain.User, error) {
	const teamQuery = `
		SELECT name
		FROM teams
		WHERE name = $1
	`

	var teamName string
	if err := r.db.QueryRowContext(ctx, teamQuery, string(name)).Scan(&teamName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Team{}, nil, repository.ErrNotFound
		}

		return domain.Team{}, nil, fmt.Errorf("get team %s: %w", name, err)
	}

	const membersQuery = `
		SELECT id, username, team_name, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, membersQuery, teamName)
	if err != nil {
		return domain.Team{}, nil, fmt.Errorf("query members for team %s: %w", name, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	members := make([]domain.User, 0)

	for rows.Next() {
		var (
			id       string
			username string
			teamName string
			isActive bool
		)

		if err := rows.Scan(&id, &username, &teamName, &isActive); err != nil {
			return domain.Team{}, nil, fmt.Errorf("scan member row for team %s: %w", name, err)
		}

		members = append(members, domain.User{
			ID:       domain.UserID(id),
			Username: username,
			TeamName: domain.TeamName(teamName),
			IsActive: isActive,
		})
	}

	if err := rows.Err(); err != nil {
		return domain.Team{}, nil, fmt.Errorf("iterate members for team %s: %w", name, err)
	}

	team := domain.Team{
		Name: name,
	}

	return team, members, nil
}

// UpsertMembers создаёт или обновляет пользователей команды по их ID.
func (r *TeamRepository) UpsertMembers(
	ctx context.Context,
	teamName domain.TeamName,
	members []domain.User,
) error {
	if len(members) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx for upsert members of team %s: %w", teamName, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const query = `
		INSERT INTO users (id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET
			username  = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active
	`

	for _, m := range members {
		_, err := tx.ExecContext(
			ctx,
			query,
			string(m.ID),
			m.Username,
			string(teamName),
			m.IsActive,
		)
		if err != nil {
			return fmt.Errorf("upsert member %s for team %s: %w", m.ID, teamName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit upsert members for team %s: %w", teamName, err)
	}

	return nil
}
