// Package postgres содержит реализацию репозиториев поверх PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/domain"
	"github.com/dixitix/pr-reviewer-service/internal/repository"
)

// PullRequestRepository реализует repository.PullRequestRepository поверх *sql.DB.
type PullRequestRepository struct {
	db *sql.DB
}

// NewPullRequestRepository создаёт новый экземпляр PullRequestRepository.
func NewPullRequestRepository(db *sql.DB) repository.PullRequestRepository {
	return &PullRequestRepository{db: db}
}

// Create создаёт новый PR и всех его ревьюверов.
func (r *PullRequestRepository) Create(ctx context.Context, pr domain.PullRequest) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const insertPR = `
		INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = tx.ExecContext(
		ctx,
		insertPR,
		pr.ID,
		pr.Name,
		pr.AuthorID,
		string(pr.Status),
		pr.CreatedAt,
		pr.MergedAt,
	)
	if err != nil {
		err = fmt.Errorf("insert pull_requests: %w", err)
		return err
	}

	const insertReviewer = `
		INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`

	for _, reviewerID := range pr.AssignedReviewers {
		if _, err = tx.ExecContext(ctx, insertReviewer, pr.ID, reviewerID); err != nil {
			err = fmt.Errorf("insert pull_request_reviewers: %w", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// GetByID возвращает PR и всех его ревьюверов.
func (r *PullRequestRepository) GetByID(
	ctx context.Context,
	id domain.PullRequestID,
) (domain.PullRequest, error) {
	const selectPR = `
		SELECT id, name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE id = $1
	`

	var (
		pr          domain.PullRequest
		statusValue string
	)

	row := r.db.QueryRowContext(ctx, selectPR, id)
	err := row.Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&statusValue,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PullRequest{}, repository.ErrNotFound
		}

		return domain.PullRequest{}, fmt.Errorf("select pull_requests: %w", err)
	}

	pr.Status = domain.PullRequestStatus(statusValue)

	const selectReviewers = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
		ORDER BY reviewer_id
	`

	rows, err := r.db.QueryContext(ctx, selectReviewers, id)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("select pull_request_reviewers: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var reviewerID domain.UserID

		if err := rows.Scan(&reviewerID); err != nil {
			return domain.PullRequest{}, fmt.Errorf("scan pull_request_reviewers: %w", err)
		}

		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
	}

	if err := rows.Err(); err != nil {
		return domain.PullRequest{}, fmt.Errorf("iterate pull_request_reviewers: %w", err)
	}

	return pr, nil
}

// Update обновляет запись pull_requests и список ревьюверов.
func (r *PullRequestRepository) Update(ctx context.Context, pr domain.PullRequest) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const updatePR = `
		UPDATE pull_requests
		SET name = $2,
		    author_id = $3,
		    status = $4,
		    merged_at = $5
		WHERE id = $1
	`

	res, err := tx.ExecContext(
		ctx,
		updatePR,
		pr.ID,
		pr.Name,
		pr.AuthorID,
		string(pr.Status),
		pr.MergedAt,
	)
	if err != nil {
		err = fmt.Errorf("update pull_requests: %w", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err = fmt.Errorf("rows affected: %w", err)
		return err
	}

	if rowsAffected == 0 {
		err = repository.ErrNotFound
		return err
	}

	const deleteReviewers = `
		DELETE FROM pull_request_reviewers
		WHERE pull_request_id = $1
	`

	if _, err = tx.ExecContext(ctx, deleteReviewers, pr.ID); err != nil {
		err = fmt.Errorf("delete pull_request_reviewers: %w", err)
		return err
	}

	const insertReviewer = `
		INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`

	for _, reviewerID := range pr.AssignedReviewers {
		if _, err = tx.ExecContext(ctx, insertReviewer, pr.ID, reviewerID); err != nil {
			err = fmt.Errorf("insert pull_request_reviewers: %w", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// ListByReviewer возвращает все PR'ы, где пользователь является ревьювером.
func (r *PullRequestRepository) ListByReviewer(
	ctx context.Context,
	reviewerID domain.UserID,
) ([]domain.PullRequest, error) {
	const query = `
		SELECT
			pr.id,
			pr.name,
			pr.author_id,
			pr.status,
			pr.created_at,
			pr.merged_at,
			r.reviewer_id
		FROM pull_requests pr
		JOIN pull_request_reviewers r
			ON r.pull_request_id = pr.id
		WHERE r.reviewer_id = $1
		ORDER BY pr.created_at DESC, pr.id, r.reviewer_id
	`

	rows, err := r.db.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("list pull_requests by reviewer: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	type prKey = domain.PullRequestID

	prByID := make(map[prKey]*domain.PullRequest)
	order := make([]prKey, 0)

	for rows.Next() {
		var (
			id          domain.PullRequestID
			name        string
			authorID    domain.UserID
			statusValue string
			createdAt   *time.Time
			mergedAt    *time.Time
			reviewer    domain.UserID
		)

		if err := rows.Scan(
			&id,
			&name,
			&authorID,
			&statusValue,
			&createdAt,
			&mergedAt,
			&reviewer,
		); err != nil {
			return nil, fmt.Errorf("scan list row: %w", err)
		}

		pr, exists := prByID[id]
		if !exists {
			newPR := &domain.PullRequest{
				ID:        id,
				Name:      name,
				AuthorID:  authorID,
				Status:    domain.PullRequestStatus(statusValue),
				CreatedAt: createdAt,
				MergedAt:  mergedAt,
			}

			prByID[id] = newPR
			order = append(order, id)
			pr = newPR
		}

		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list rows: %w", err)
	}

	result := make([]domain.PullRequest, 0, len(order))
	for _, id := range order {
		pr := prByID[id]
		result = append(result, *pr)
	}

	return result, nil
}
