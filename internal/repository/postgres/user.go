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

// UserRepository реализует repository.UserRepository поверх *sql.DB.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создаёт новый экземпляр UserRepository.
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// GetByID возвращает пользователя по его идентификатору.
func (r *UserRepository) GetByID(
	ctx context.Context,
	id domain.UserID,
) (domain.User, error) {
	const query = `
		SELECT id, username, team_name, is_active
		FROM users
		WHERE id = $1
	`

	var (
		userID   string
		username string
		teamName string
		isActive bool
	)

	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(&userID, &username, &teamName, &isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, repository.ErrNotFound
		}

		return domain.User{}, fmt.Errorf("get user %s: %w", id, err)
	}

	return domain.User{
		ID:       domain.UserID(userID),
		Username: username,
		TeamName: domain.TeamName(teamName),
		IsActive: isActive,
	}, nil
}

// SetActive меняет флаг активности пользователя.
func (r *UserRepository) SetActive(
	ctx context.Context,
	id domain.UserID,
	isActive bool,
) error {
	const query = `
		UPDATE users
		SET is_active = $2
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, string(id), isActive)
	if err != nil {
		return fmt.Errorf("update user %s active=%t: %w", id, isActive, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for user %s: %w", id, err)
	}

	if rows == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// ListActiveByTeam возвращает активных пользователей команды.
// Если excludeID != nil, этот пользователь исключается из результата.
func (r *UserRepository) ListActiveByTeam(
	ctx context.Context,
	teamName domain.TeamName,
	excludeID *domain.UserID,
) ([]domain.User, error) {
	baseQuery := `
		SELECT id, username, team_name, is_active
		FROM users
		WHERE team_name = $1
		  AND is_active = TRUE
	`

	args := []any{string(teamName)}

	// Динамически добавляем фильтр исключения, если нужно.
	if excludeID != nil {
		baseQuery += " AND id <> $2"
		args = append(args, string(*excludeID))
	}

	baseQuery += " ORDER BY id"

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("list active users for team %s: %w", teamName, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	users := make([]domain.User, 0)

	for rows.Next() {
		var (
			id       string
			username string
			tName    string
			isActive bool
		)

		if err := rows.Scan(&id, &username, &tName, &isActive); err != nil {
			return nil, fmt.Errorf("scan active user for team %s: %w", teamName, err)
		}

		users = append(users, domain.User{
			ID:       domain.UserID(id),
			Username: username,
			TeamName: domain.TeamName(tName),
			IsActive: isActive,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active users for team %s: %w", teamName, err)
	}

	if users == nil {
		return []domain.User{}, nil
	}

	return users, nil
}
