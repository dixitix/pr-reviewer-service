// Package postgres_test содержит интеграционные тесты репозиториев.
package integration

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
	"github.com/dixitix/pr-reviewer-service/internal/repository/postgres"
)

// newTestUserRepository создаёт UserRepository поверх тестовой БД.
func newTestUserRepository(t *testing.T) (*sql.DB, repository.UserRepository) {
	db := openTestDB(t)
	truncateAllTables(t, db)

	return db, postgres.NewUserRepository(db)
}

// TestUserRepository_GetByID проверяет получение пользователя по идентификатору.
func TestUserRepository_GetByID(t *testing.T) {
	db, repo := newTestUserRepository(t)

	const (
		teamName = "backend"
		userID   = "user-1"
	)

	insertTeam(t, db, teamName)
	insertUser(t, db, userID, "alice", teamName, true)

	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()

		user, err := repo.GetByID(ctx, domain.UserID(userID))
		if err != nil {
			t.Fatalf("GetByID returned error: %v", err)
		}

		if user.ID != domain.UserID(userID) {
			t.Fatalf("unexpected user ID: got %s, want %s", user.ID, userID)
		}

		if user.Username != "alice" {
			t.Fatalf("unexpected username: got %s, want %s", user.Username, "alice")
		}

		if user.TeamName != domain.TeamName(teamName) {
			t.Fatalf("unexpected team name: got %s, want %s", user.TeamName, teamName)
		}

		if !user.IsActive {
			t.Fatalf("expected user to be active")
		}
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()

		_, err := repo.GetByID(ctx, domain.UserID("unknown"))

		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got: %v", err)
		}
	})
}

// TestUserRepository_SetActive проверяет изменение статуса активности пользователя.
func TestUserRepository_SetActive(t *testing.T) {
	db, repo := newTestUserRepository(t)

	const (
		teamName = "backend"
		userID   = "user-2"
	)

	insertTeam(t, db, teamName)
	insertUser(t, db, userID, "bob", teamName, false)

	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()

		if err := repo.SetActive(ctx, domain.UserID(userID), true); err != nil {
			t.Fatalf("SetActive returned error: %v", err)
		}

		user, err := repo.GetByID(ctx, domain.UserID(userID))
		if err != nil {
			t.Fatalf("GetByID after SetActive returned error: %v", err)
		}

		if !user.IsActive {
			t.Fatalf("expected user to be active after SetActive")
		}
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()

		err := repo.SetActive(ctx, domain.UserID("unknown"), true)
		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got: %v", err)
		}
	})
}

// TestUserRepository_ListActiveByTeam проверяет выборку активных пользователей команды.
func TestUserRepository_ListActiveByTeam(t *testing.T) {
	db, repo := newTestUserRepository(t)

	const (
		teamA = "team-a"
		teamB = "team-b"
	)

	insertTeam(t, db, teamA)
	insertTeam(t, db, teamB)

	insertUser(t, db, "user-1", "alice", teamA, true)
	insertUser(t, db, "user-2", "bob", teamA, true)
	insertUser(t, db, "user-3", "charlie", teamA, false)
	insertUser(t, db, "user-4", "dave", teamB, true)

	ctx := context.Background()

	t.Run("only active users of team", func(t *testing.T) {
		users, err := repo.ListActiveByTeam(ctx, domain.TeamName(teamA), nil)
		if err != nil {
			t.Fatalf("ListActiveByTeam returned error: %v", err)
		}

		if len(users) != 2 {
			t.Fatalf("unexpected number of users: got %d, want %d", len(users), 2)
		}

		if users[0].ID != domain.UserID("user-1") || users[1].ID != domain.UserID("user-2") {
			t.Fatalf("unexpected users order or ids: got [%s, %s]", users[0].ID, users[1].ID)
		}
	})

	t.Run("exclude specific user", func(t *testing.T) {
		exclude := domain.UserID("user-1")

		users, err := repo.ListActiveByTeam(ctx, domain.TeamName(teamA), &exclude)
		if err != nil {
			t.Fatalf("ListActiveByTeam with exclude returned error: %v", err)
		}

		if len(users) != 1 {
			t.Fatalf("unexpected number of users: got %d, want %d", len(users), 1)
		}

		if users[0].ID != domain.UserID("user-2") {
			t.Fatalf("unexpected user id: got %s, want %s", users[0].ID, "user-2")
		}
	})

	t.Run("empty result", func(t *testing.T) {
		users, err := repo.ListActiveByTeam(ctx, domain.TeamName(teamB), nil)
		if err != nil {
			t.Fatalf("ListActiveByTeam for teamB returned error: %v", err)
		}

		if len(users) != 1 {
			t.Fatalf("unexpected number of users for teamB: got %d, want %d", len(users), 1)
		}

		if users[0].ID != domain.UserID("user-4") {
			t.Fatalf("unexpected user for teamB: got %s, want %s", users[0].ID, "user-4")
		}
	})
}
