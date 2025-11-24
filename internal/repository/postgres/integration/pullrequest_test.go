// Package postgres_test содержит интеграционные тесты репозиториев.
package integration

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
	"github.com/dixitix/pr-reviewer-service/internal/repository/postgres"
)

// newTestPullRequestRepository создаёт PullRequestRepository поверх тестовой БД.
func newTestPullRequestRepository(t *testing.T) (*sql.DB, repository.PullRequestRepository) {
	db := openTestDB(t)
	truncateAllTables(t, db)

	return db, postgres.NewPullRequestRepository(db)
}

// TestPullRequestRepository_CreateAndGet проверяет успешное создание и получение pull request.
func TestPullRequestRepository_CreateAndGet(t *testing.T) {
	db, repo := newTestPullRequestRepository(t)
	defer db.Close()

	ctx := context.Background()

	const (
		teamName    = "backend"
		authorID    = "author-1"
		reviewer1   = "reviewer-1"
		reviewer2   = "reviewer-2"
		pullReqID   = "pr-1"
		pullReqName = "Add feature"
	)

	insertTeam(t, db, teamName)
	insertUser(t, db, authorID, "author", teamName, true)
	insertUser(t, db, reviewer1, "r1", teamName, true)
	insertUser(t, db, reviewer2, "r2", teamName, true)

	now := time.Now().UTC().Truncate(time.Second)

	pr := domain.PullRequest{
		ID:                domain.PullRequestID(pullReqID),
		Name:              pullReqName,
		AuthorID:          domain.UserID(authorID),
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: []domain.UserID{domain.UserID(reviewer1), domain.UserID(reviewer2)},
		CreatedAt:         &now,
		// MergedAt оставляем nil.
	}

	if err := repo.Create(ctx, pr); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	got, err := repo.GetByID(ctx, pr.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if got.ID != pr.ID {
		t.Fatalf("ID mismatch: got %q, want %q", got.ID, pr.ID)
	}
	if got.Name != pr.Name {
		t.Fatalf("Name mismatch: got %q, want %q", got.Name, pr.Name)
	}
	if got.AuthorID != pr.AuthorID {
		t.Fatalf("AuthorID mismatch: got %q, want %q", got.AuthorID, pr.AuthorID)
	}
	if got.Status != pr.Status {
		t.Fatalf("Status mismatch: got %q, want %q", got.Status, pr.Status)
	}
	if got.CreatedAt == nil {
		t.Fatalf("CreatedAt is nil, expected non-nil")
	}

	if len(got.AssignedReviewers) != len(pr.AssignedReviewers) {
		t.Fatalf("AssignedReviewers length mismatch: got %d, want %d", len(got.AssignedReviewers), len(pr.AssignedReviewers))
	}

	for i := range pr.AssignedReviewers {
		if got.AssignedReviewers[i] != pr.AssignedReviewers[i] {
			t.Fatalf("AssignedReviewers[%d] mismatch: got %q, want %q",
				i, got.AssignedReviewers[i], pr.AssignedReviewers[i])
		}
	}
}

// TestPullRequestRepository_Update убеждается, что обновление полей pull request сохраняется.
func TestPullRequestRepository_Update(t *testing.T) {
	db, repo := newTestPullRequestRepository(t)
	defer db.Close()

	ctx := context.Background()

	const (
		teamName    = "backend"
		authorID    = "author-2"
		reviewer1   = "reviewer-3"
		reviewer2   = "reviewer-4"
		pullReqID   = "pr-2"
		pullReqName = "Initial name"
	)

	insertTeam(t, db, teamName)
	insertUser(t, db, authorID, "author2", teamName, true)
	insertUser(t, db, reviewer1, "r3", teamName, true)
	insertUser(t, db, reviewer2, "r4", teamName, true)

	createdAt := time.Now().UTC().Truncate(time.Second)

	pr := domain.PullRequest{
		ID:                domain.PullRequestID(pullReqID),
		Name:              pullReqName,
		AuthorID:          domain.UserID(authorID),
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: []domain.UserID{domain.UserID(reviewer1), domain.UserID(reviewer2)},
		CreatedAt:         &createdAt,
	}

	if err := repo.Create(ctx, pr); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	mergedAt := time.Now().UTC().Truncate(time.Second)
	pr.Status = domain.PullRequestStatusMerged
	pr.MergedAt = &mergedAt
	pr.AssignedReviewers = []domain.UserID{domain.UserID(reviewer2)}

	if err := repo.Update(ctx, pr); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	got, err := repo.GetByID(ctx, pr.ID)
	if err != nil {
		t.Fatalf("GetByID after Update returned error: %v", err)
	}

	if got.Status != domain.PullRequestStatusMerged {
		t.Fatalf("Status mismatch after update: got %q, want %q", got.Status, domain.PullRequestStatusMerged)
	}
	if got.MergedAt == nil {
		t.Fatalf("MergedAt is nil after update, expected non-nil")
	}
	if len(got.AssignedReviewers) != 1 {
		t.Fatalf("AssignedReviewers length after update: got %d, want %d", len(got.AssignedReviewers), 1)
	}
	if got.AssignedReviewers[0] != domain.UserID(reviewer2) {
		t.Fatalf("AssignedReviewers[0] mismatch: got %q, want %q", got.AssignedReviewers[0], reviewer2)
	}
}

// TestPullRequestRepository_Update_NotFound проверяет, что обновление несуществующего pull request возвращает ErrNotFound.
func TestPullRequestRepository_Update_NotFound(t *testing.T) {
	db, repo := newTestPullRequestRepository(t)
	defer db.Close()

	ctx := context.Background()

	pr := domain.PullRequest{
		ID:       domain.PullRequestID("non-existent"),
		Name:     "Does not matter",
		AuthorID: domain.UserID("unknown"),
		Status:   domain.PullRequestStatusOpen,
	}

	err := repo.Update(ctx, pr)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}

// TestPullRequestRepository_ListByReviewer убеждается, что репозиторий возвращает все pull request для ревьюера.
func TestPullRequestRepository_ListByReviewer(t *testing.T) {
	db, repo := newTestPullRequestRepository(t)
	defer db.Close()

	ctx := context.Background()

	const (
		teamName   = "backend"
		authorID   = "author-3"
		reviewerID = "reviewer-5"
		otherRevID = "reviewer-6"
	)

	insertTeam(t, db, teamName)
	insertUser(t, db, authorID, "author3", teamName, true)
	insertUser(t, db, reviewerID, "rev5", teamName, true)
	insertUser(t, db, otherRevID, "rev6", teamName, true)

	now := time.Now().UTC().Truncate(time.Second)
	now2 := now.Add(1 * time.Minute)

	pr1 := domain.PullRequest{
		ID:                domain.PullRequestID("pr-3"),
		Name:              "First PR",
		AuthorID:          domain.UserID(authorID),
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: []domain.UserID{domain.UserID(reviewerID), domain.UserID(otherRevID)},
		CreatedAt:         &now,
	}
	pr2 := domain.PullRequest{
		ID:                domain.PullRequestID("pr-4"),
		Name:              "Second PR",
		AuthorID:          domain.UserID(authorID),
		Status:            domain.PullRequestStatusMerged,
		AssignedReviewers: []domain.UserID{domain.UserID(reviewerID)},
		CreatedAt:         &now2,
	}

	if err := repo.Create(ctx, pr1); err != nil {
		t.Fatalf("Create pr1 returned error: %v", err)
	}
	if err := repo.Create(ctx, pr2); err != nil {
		t.Fatalf("Create pr2 returned error: %v", err)
	}

	list, err := repo.ListByReviewer(ctx, domain.UserID(reviewerID))
	if err != nil {
		t.Fatalf("ListByReviewer returned error: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("unexpected number of PRs: got %d, want %d", len(list), 2)
	}

	found := make(map[domain.PullRequestID]domain.PullRequest)
	for _, pr := range list {
		found[pr.ID] = pr
	}

	got1, ok := found[pr1.ID]
	if !ok {
		t.Fatalf("PR %q not found in result", pr1.ID)
	}
	got2, ok := found[pr2.ID]
	if !ok {
		t.Fatalf("PR %q not found in result", pr2.ID)
	}

	if len(got1.AssignedReviewers) != 1 || got1.AssignedReviewers[0] != domain.UserID(reviewerID) {
		t.Fatalf("unexpected reviewers in pr1: %+v", got1.AssignedReviewers)
	}
	if len(got2.AssignedReviewers) != 1 || got2.AssignedReviewers[0] != domain.UserID(reviewerID) {
		t.Fatalf("unexpected reviewers in pr2: %+v", got2.AssignedReviewers)
	}
}
