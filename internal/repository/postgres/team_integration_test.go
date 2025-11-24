//go:build integration

// Package postgres_test содержит интеграционные тесты репозиториев.
package postgres_test

import (
	"context"
	"database/sql"
	"sort"
	"testing"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
	"github.com/dixitix/pr-reviewer-service/internal/repository/postgres"
)

// newTeamRepo подготавливает чистую тестовую базу и репозиторий команд.
func newTeamRepo(t *testing.T) (*sql.DB, repository.TeamRepository) {
	t.Helper()

	db := openTestDB(t)
	truncateAllTables(t, db)

	repo := postgres.NewTeamRepository(db)

	return db, repo
}

// TestTeamRepository_CreateAndGet проверяет создание команды и получение ее участников.
func TestTeamRepository_CreateAndGet(t *testing.T) {
	_, repo := newTeamRepo(t)

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
