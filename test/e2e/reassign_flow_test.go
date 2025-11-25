// Package e2e содержит e2e тесты.
package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/http/httperr"
	"github.com/dixitix/pr-reviewer-service/internal/http/pullrequest"
	"github.com/dixitix/pr-reviewer-service/internal/http/team"
)

// TestE2E_ReassignFlow:
// 1) /team/add
// 2) /pullRequest/create
// 3) /pullRequest/reassign
// 4) /pullRequest/merge
// 5) /pullRequest/reassign ещё раз => 409 + PR_MERGED
func TestE2E_ReassignFlow(t *testing.T) {
	suffix := time.Now().UnixNano()

	teamName := fmt.Sprintf("backend-e2e-reassign-%d", suffix)
	authorID := fmt.Sprintf("u-author-%d", suffix)
	r1ID := fmt.Sprintf("u-r1-%d", suffix)
	r2ID := fmt.Sprintf("u-r2-%d", suffix)
	r3ID := fmt.Sprintf("u-r3-%d", suffix)

	prID := fmt.Sprintf("pr-reassign-%d", suffix)
	prName := "Reassign flow e2e"

	teamReq := team.DTO{
		TeamName: teamName,
		Members: []team.MemberDTO{
			{UserID: authorID, Username: "Author", IsActive: true},
			{UserID: r1ID, Username: "Reviewer1", IsActive: true},
			{UserID: r2ID, Username: "Reviewer2", IsActive: true},
			{UserID: r3ID, Username: "Reviewer3", IsActive: true},
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

	if len(pr.AssignedReviewers) == 0 {
		t.Fatalf("no reviewers assigned on create")
	}

	oldReviewer := pr.AssignedReviewers[0]

	reassignReq := pullrequest.ReassignPullRequestRequest{
		PullRequestID: prID,
		OldUserID:     oldReviewer,
	}

	var reassignResp pullrequest.ReassignResponse
	doRequest(
		t,
		http.MethodPost,
		"/pullRequest/reassign",
		reassignReq,
		http.StatusOK,
		&reassignResp,
	)

	if reassignResp.PullRequest.PullRequestID != prID {
		t.Fatalf("reassign: pull_request_id mismatch: got %q, want %q", reassignResp.PullRequest.PullRequestID, prID)
	}
	if reassignResp.ReplacedBy == "" {
		t.Fatalf("reassign: replaced_by is empty")
	}
	if reassignResp.ReplacedBy == oldReviewer {
		t.Fatalf("reassign: replaced_by == oldReviewer (%s)", oldReviewer)
	}
	if reassignResp.ReplacedBy == authorID {
		t.Fatalf("reassign: replaced_by must not be author (%s)", authorID)
	}

	for _, id := range reassignResp.PullRequest.AssignedReviewers {
		if id == oldReviewer {
			t.Fatalf("reassign: old reviewer %s still in assigned_reviewers: %+v", oldReviewer, reassignResp.PullRequest.AssignedReviewers)
		}
	}

	foundNew := false
	for _, id := range reassignResp.PullRequest.AssignedReviewers {
		if id == reassignResp.ReplacedBy {
			foundNew = true
		}
	}
	if !foundNew {
		t.Fatalf("reassign: replaced_by %s is not in assigned_reviewers: %+v", reassignResp.ReplacedBy, reassignResp.PullRequest.AssignedReviewers)
	}

	mergeReq := pullrequest.MergePullRequestRequest{PullRequestID: prID}

	var mergeResp pullrequest.Envelope
	doRequest(
		t,
		http.MethodPost,
		"/pullRequest/merge",
		mergeReq,
		http.StatusOK,
		&mergeResp,
	)

	if mergeResp.PullRequest.Status != "MERGED" {
		t.Fatalf("status after merge: got %q, want %q", mergeResp.PullRequest.Status, "MERGED")
	}

	reassignReq2 := pullrequest.ReassignPullRequestRequest{
		PullRequestID: prID,
		OldUserID:     reassignResp.ReplacedBy,
	}

	var errResp httperr.ErrorResponse
	doRequest(
		t,
		http.MethodPost,
		"/pullRequest/reassign",
		reassignReq2,
		http.StatusConflict,
		&errResp,
	)

	if errResp.Error.Code != "PR_MERGED" {
		t.Fatalf("expected error code PR_MERGED, got %q (%s)", errResp.Error.Code, errResp.Error.Message)
	}
}
