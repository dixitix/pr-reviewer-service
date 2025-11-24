//go:build integration

// Package postgres_test содержит интеграционные тесты репозиториев.
package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"testing"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
	"github.com/dixitix/pr-reviewer-service/internal/repository/postgres"
)

// newTeamRepo подготавливает чистую тестовую базу и репозиторий команд.
func newTestTeamRepository(t *testing.T) (*sql.DB, repository.TeamRepository) {
	db := openTestDB(t)
	truncateAllTables(t, db)

	return db, postgres.NewTeamRepository(db)
}

// TestTeamRepository_CreateAndGet проверяет создание команды и получение ее участников.
func TestTeamRepository_CreateAndGet(t *testing.T) {
	_, repo := newTestTeamRepository(t)

	ctx := context.Background()

	const teamName = domain.TeamName("backend")

	team := domain.Team{
		Name: teamName,
	}

	members := []domain.User{
		{
			ID:       domain.UserID("u1"),
			Username: "Alice",
			TeamName: teamName,
			IsActive: true,
		},
		{
			ID:       domain.UserID("u2"),
			Username: "Bob",
			TeamName: teamName,
			IsActive: true,
		},
	}

	if err := repo.CreateTeam(ctx, team); err != nil {
		t.Fatalf("CreateTeam: %v", err)
	}

	if err := repo.UpsertMembers(ctx, teamName, members); err != nil {
		t.Fatalf("UpsertMembers: %v", err)
	}

	gotTeam, gotMembers, err := repo.GetTeamWithMembers(ctx, teamName)
	if err != nil {
		t.Fatalf("GetTeamWithMembers: %v", err)
	}

	if gotTeam.Name != teamName {
		t.Fatalf("team name mismatch: got %q, want %q", gotTeam.Name, teamName)
	}

	if len(gotMembers) != len(members) {
		t.Fatalf("members len mismatch: got %d, want %d", len(gotMembers), len(members))
	}

	sort.Slice(gotMembers, func(i, j int) bool {
		return string(gotMembers[i].ID) < string(gotMembers[j].ID)
	})
	sort.Slice(members, func(i, j int) bool {
		return string(members[i].ID) < string(members[j].ID)
	})

	for i := range members {
		if gotMembers[i].ID != members[i].ID {
			t.Fatalf("member[%d].ID mismatch: got %q, want %q", i, gotMembers[i].ID, members[i].ID)
		}
		if gotMembers[i].Username != members[i].Username {
			t.Fatalf("member[%d].Username mismatch: got %q, want %q", i, gotMembers[i].Username, members[i].Username)
		}
		if gotMembers[i].TeamName != members[i].TeamName {
			t.Fatalf("member[%d].TeamName mismatch: got %q, want %q", i, gotMembers[i].TeamName, members[i].TeamName)
		}
		if gotMembers[i].IsActive != members[i].IsActive {
			t.Fatalf("member[%d].IsActive mismatch: got %v, want %v", i, gotMembers[i].IsActive, members[i].IsActive)
		}
	}
}

// TestTeamRepository_TeamExists проверяет, что TeamExists корректно работает для существующей и несуществующей команды.
func TestTeamRepository_TeamExists(t *testing.T) {
	_, repo := newTestTeamRepository(t)
	ctx := context.Background()

	const teamName = domain.TeamName("backend")

	exists, err := repo.TeamExists(ctx, teamName)
	if err != nil {
		t.Fatalf("TeamExists(before CreateTeam): %v", err)
	}
	if exists {
		t.Fatalf("TeamExists(before CreateTeam): got %v, want false", exists)
	}

	if err := repo.CreateTeam(ctx, domain.Team{Name: teamName}); err != nil {
		t.Fatalf("CreateTeam: %v", err)
	}

	exists, err = repo.TeamExists(ctx, teamName)
	if err != nil {
		t.Fatalf("TeamExists(after CreateTeam): %v", err)
	}
	if !exists {
		t.Fatalf("TeamExists(after CreateTeam): got %v, want true", exists)
	}
}

// TestTeamRepository_GetTeamWithMembers_NotFound проверяет, что для несуществующей команды возвращается repository.ErrNotFound.
func TestTeamRepository_GetTeamWithMembers_NotFound(t *testing.T) {
	_, repo := newTestTeamRepository(t)
	ctx := context.Background()

	_, _, err := repo.GetTeamWithMembers(ctx, domain.TeamName("unknown"))
	if err == nil {
		t.Fatal("GetTeamWithMembers: got nil error, want ErrNotFound")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("GetTeamWithMembers: got %v, want ErrNotFound", err)
	}
}

// TestTeamRepository_UpsertMembers_UpdateAndInsert проверяет, что UpsertMembers обновляет существующих участников и добавляет новых.
func TestTeamRepository_UpsertMembers_UpdateAndInsert(t *testing.T) {
	_, repo := newTestTeamRepository(t)
	ctx := context.Background()

	const teamName = domain.TeamName("backend")

	team := domain.Team{Name: teamName}
	if err := repo.CreateTeam(ctx, team); err != nil {
		t.Fatalf("CreateTeam: %v", err)
	}

	initialMembers := []domain.User{
		{
			ID:       domain.UserID("u1"),
			Username: "Alice",
			TeamName: teamName,
			IsActive: true,
		},
		{
			ID:       domain.UserID("u2"),
			Username: "Bob",
			TeamName: teamName,
			IsActive: true,
		},
	}

	if err := repo.UpsertMembers(ctx, teamName, initialMembers); err != nil {
		t.Fatalf("UpsertMembers(initial): %v", err)
	}

	updatedMembers := []domain.User{
		{
			ID:       domain.UserID("u1"),
			Username: "Alice Renamed",
			TeamName: teamName,
			IsActive: false,
		},
		{
			ID:       domain.UserID("u2"),
			Username: "Bob",
			TeamName: teamName,
			IsActive: true,
		},
		{
			ID:       domain.UserID("u3"),
			Username: "Charlie",
			TeamName: teamName,
			IsActive: true,
		},
	}

	if err := repo.UpsertMembers(ctx, teamName, updatedMembers); err != nil {
		t.Fatalf("UpsertMembers(updated): %v", err)
	}

	gotTeam, gotMembers, err := repo.GetTeamWithMembers(ctx, teamName)
	if err != nil {
		t.Fatalf("GetTeamWithMembers: %v", err)
	}

	if gotTeam.Name != teamName {
		t.Fatalf("team name mismatch: got %q, want %q", gotTeam.Name, teamName)
	}

	if len(gotMembers) != len(updatedMembers) {
		t.Fatalf("members len mismatch: got %d, want %d", len(gotMembers), len(updatedMembers))
	}

	sort.Slice(gotMembers, func(i, j int) bool {
		return string(gotMembers[i].ID) < string(gotMembers[j].ID)
	})
	sort.Slice(updatedMembers, func(i, j int) bool {
		return string(updatedMembers[i].ID) < string(updatedMembers[j].ID)
	})

	for i := range updatedMembers {
		if gotMembers[i].ID != updatedMembers[i].ID {
			t.Fatalf("member[%d].ID mismatch: got %q, want %q", i, gotMembers[i].ID, updatedMembers[i].ID)
		}
		if gotMembers[i].Username != updatedMembers[i].Username {
			t.Fatalf("member[%d].Username mismatch: got %q, want %q", i, gotMembers[i].Username, updatedMembers[i].Username)
		}
		if gotMembers[i].TeamName != updatedMembers[i].TeamName {
			t.Fatalf("member[%d].TeamName mismatch: got %q, want %q", i, gotMembers[i].TeamName, updatedMembers[i].TeamName)
		}
		if gotMembers[i].IsActive != updatedMembers[i].IsActive {
			t.Fatalf("member[%d].IsActive mismatch: got %v, want %v", i, gotMembers[i].IsActive, updatedMembers[i].IsActive)
		}
	}
}
