// Package e2e содержит e2e тесты.
package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/http/pullrequest"
	"github.com/dixitix/pr-reviewer-service/internal/http/team"
	"github.com/dixitix/pr-reviewer-service/internal/http/user"
)

// TestE2E_CreateMergeAndReview:
// 1) /team/add
// 2) /pullRequest/create
// 3) /pullRequest/merge (2 раза — идемпотентность)
// 4) /users/getReview
func TestE2E_CreateMergeAndReview(t *testing.T) {
	suffix := time.Now().UnixNano()

	teamName := fmt.Sprintf("backend-e2e-%d", suffix)
	authorID := fmt.Sprintf("u-author-%d", suffix)
	reviewerID := fmt.Sprintf("u-reviewer-%d", suffix)
	otherID := fmt.Sprintf("u-other-%d", suffix)

	prID := fmt.Sprintf("pr-%d", suffix)
	prName := "Add search via e2e"

	teamReq := team.DTO{
		TeamName: teamName,
		Members: []team.MemberDTO{
			{
				UserID:   authorID,
				Username: "Author",
				IsActive: true,
			},
			{
				UserID:   reviewerID,
				Username: "Reviewer1",
				IsActive: true,
			},
			{
				UserID:   otherID,
				Username: "Reviewer2",
				IsActive: true,
			},
		},
	}

	var teamResp team.GetTeamResponse
	doRequest(
		t,
		http.MethodPost,
		"/team/add",
		teamReq,
		http.StatusCreated,
		&teamResp,
	)

	if teamResp.Team.TeamName != teamName {
		t.Fatalf("team_name mismatch: got %q, want %q", teamResp.Team.TeamName, teamName)
	}
	if len(teamResp.Team.Members) != 3 {
		t.Fatalf("expected 3 team members, got %d", len(teamResp.Team.Members))
	}

	createReq := pullrequest.CreatePullRequestRequest{
		PullRequestID:   prID,
		PullRequestName: prName,
		AuthorID:        authorID,
	}

	var createResp pullrequest.Envelope
	doRequest(
		t,
		http.MethodPost,
		"/pullRequest/create",
		createReq,
		http.StatusCreated,
		&createResp,
	)

	pr := createResp.PullRequest

	if pr.PullRequestID != prID {
		t.Fatalf("pull_request_id mismatch: got %q, want %q", pr.PullRequestID, prID)
	}
	if pr.Status != "OPEN" {
		t.Fatalf("status after create: got %q, want %q", pr.Status, "OPEN")
	}
	if len(pr.AssignedReviewers) == 0 || len(pr.AssignedReviewers) > 2 {
		t.Fatalf("assigned_reviewers length = %d, want 1..2", len(pr.AssignedReviewers))
	}
	for _, id := range pr.AssignedReviewers {
		if id == authorID {
			t.Fatalf("author %q must not be in assigned_reviewers: %+v", authorID, pr.AssignedReviewers)
		}
	}

	mergeReq := pullrequest.MergePullRequestRequest{PullRequestID: prID}

	var mergeResp1 pullrequest.Envelope
	doRequest(
		t,
		http.MethodPost,
		"/pullRequest/merge",
		mergeReq,
		http.StatusOK,
		&mergeResp1,
	)

	if mergeResp1.PullRequest.Status != "MERGED" {
		t.Fatalf("status after first merge: got %q, want %q", mergeResp1.PullRequest.Status, "MERGED")
	}
	if mergeResp1.PullRequest.MergedAt == nil {
		t.Fatalf("MergedAt is nil after first merge, want non-nil")
	}

	var mergeResp2 pullrequest.Envelope
	doRequest(
		t,
		http.MethodPost,
		"/pullRequest/merge",
		mergeReq,
		http.StatusOK,
		&mergeResp2,
	)

	if mergeResp2.PullRequest.Status != "MERGED" {
		t.Fatalf("status after second merge: got %q, want %q", mergeResp2.PullRequest.Status, "MERGED")
	}

	var reviewResp user.GetUserReviewResponse
	path := fmt.Sprintf("/users/getReview?user_id=%s", reviewerID)
	doRequest(
		t,
		http.MethodGet,
		path,
		nil,
		http.StatusOK,
		&reviewResp,
	)

	if reviewResp.UserID != reviewerID {
		t.Fatalf("user_id in getReview: got %q, want %q", reviewResp.UserID, reviewerID)
	}

	found := false
	for _, short := range reviewResp.PullRequests {
		if short.PullRequestID == prID {
			found = true
			if short.Status != "MERGED" {
				t.Fatalf("PR %s status in getReview: got %q, want %q", prID, short.Status, "MERGED")
			}
		}
	}
	if !found {
		t.Fatalf("PR %s not found in getReview response", prID)
	}
}
