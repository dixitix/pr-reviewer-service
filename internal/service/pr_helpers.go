// Package service содержит реализацию бизнес-логики сервиса назначения ревьюеров.
package service

import "github.com/dixitix/pr-reviewer-service/internal/domain"

// pickReviewersForNewPR выбирает до limit ревьюверов из списка активных участников команды.
func (s *service) pickReviewersForNewPR(
	users []domain.User,
	limit int,
) []domain.UserID {
	if len(users) == 0 || limit <= 0 {
		return nil
	}

	ids := make([]domain.UserID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}

	s.rndMu.Lock()
	s.rnd.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	s.rndMu.Unlock()

	if len(ids) <= limit {
		return ids
	}

	return ids[:limit]
}

// buildReplacementCandidates подбирает кандидатов для замены ревьювера.
func (s *service) buildReplacementCandidates(
	pr domain.PullRequest,
	reviewerID domain.UserID,
	activeMembers []domain.User,
) []domain.UserID {
	if len(activeMembers) == 0 {
		return nil
	}

	assigned := make(map[domain.UserID]struct{}, len(pr.AssignedReviewers))
	for _, id := range pr.AssignedReviewers {
		assigned[id] = struct{}{}
	}

	candidates := make([]domain.UserID, 0, len(activeMembers))
	for _, u := range activeMembers {
		// Не назначаем автора.
		if u.ID == pr.AuthorID {
			continue
		}

		// Не оставляем заменяемого ревьювера.
		if u.ID == reviewerID {
			continue
		}

		// Не назначаем уже присутствующих ревьюверов.
		if _, exists := assigned[u.ID]; exists {
			continue
		}

		candidates = append(candidates, u.ID)
	}

	return candidates
}
