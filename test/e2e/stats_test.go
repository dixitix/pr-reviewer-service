// Package e2e содержит e2e тесты.
package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/http/pullrequest"
	"github.com/dixitix/pr-reviewer-service/internal/http/team"
)

// TestE2E_AssignmentsStats проверяет, что статистика по пользователям и PR возвращает корректные данные.
func TestE2E_AssignmentsStats(t *testing.T) {
	suffix := time.Now().UnixNano()

	teamName := fmt.Sprintf("stats-team-%d", suffix)
	authorID := fmt.Sprintf("u-author-%d", suffix)
	r1ID := fmt.Sprintf("u-r1-%d", suffix)
	r2ID := fmt.Sprintf("u-r2-%d", suffix)

	teamReq := team.DTO{
		TeamName: teamName,
		Members: []team.MemberDTO{
			{UserID: authorID, Username: "Author", IsActive: true},
			{UserID: r1ID, Username: "Reviewer1", IsActive: true},
			{UserID: r2ID, Username: "Reviewer2", IsActive: true},
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

	prIDs := []string{
		fmt.Sprintf("pr-stats-%d-1", suffix),
		fmt.Sprintf("pr-stats-%d-2", suffix),
	}

	expectedUserAssignments := map[string]int{
		r1ID: 0,
		r2ID: 0,
	}
	expectedPRAssignments := make(map[string]int, len(prIDs))

	for _, prID := range prIDs {
		req := pullrequest.CreatePullRequestRequest{
			PullRequestID:   prID,
			PullRequestName: "Stats PR",
			AuthorID:        authorID,
		}

		var resp pullrequest.Envelope
		doRequest(
			t,
			http.MethodPost,
			"/pullRequest/create",
			req,
			http.StatusCreated,
			&resp,
		)

		assigned := resp.PullRequest.AssignedReviewers
		if len(assigned) == 0 {
			t.Fatalf("PR %s has no assigned reviewers", prID)
		}

		expectedPRAssignments[prID] = len(assigned)
		for _, reviewer := range assigned {
			if reviewer == authorID {
				t.Fatalf("author %s must not be assigned as reviewer", authorID)
			}
			expectedUserAssignments[reviewer]++
		}
	}

	var userStats struct {
		Stats []struct {
			UserID      string `json:"user_id"`
			Assignments int    `json:"assignments"`
		} `json:"stats"`
	}

	doRequest(
		t,
		http.MethodGet,
		"/stats/byUser",
		nil,
		http.StatusOK,
		&userStats,
	)

	for userID, want := range expectedUserAssignments {
		found := false
		for _, stat := range userStats.Stats {
			if stat.UserID == userID {
				found = true
				if stat.Assignments != want {
					t.Fatalf("assignments for user %s: got %d, want %d", userID, stat.Assignments, want)
				}
			}
		}
		if !found {
			t.Fatalf("stat for user %s not found in response", userID)
		}
	}

	var prStats struct {
		Stats []struct {
			PullRequestID string `json:"pull_request_id"`
			Assignments   int    `json:"assignments"`
		} `json:"stats"`
	}

	doRequest(
		t,
		http.MethodGet,
		"/stats/byPullRequest",
		nil,
		http.StatusOK,
		&prStats,
	)

	for prID, want := range expectedPRAssignments {
		found := false
		for _, stat := range prStats.Stats {
			if stat.PullRequestID == prID {
				found = true
				if stat.Assignments != want {
					t.Fatalf("assignments for PR %s: got %d, want %d", prID, stat.Assignments, want)
				}
			}
		}
		if !found {
			t.Fatalf("stat for PR %s not found in response", prID)
		}
	}
}
