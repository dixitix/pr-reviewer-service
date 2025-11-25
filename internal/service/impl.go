// Package service содержит реализацию бизнес-логики сервиса назначения ревьюеров.
package service

import (
	"math/rand"
	"sync"
	"time"

	"github.com/dixitix/pr-reviewer-service/internal/repository"
)

// service — реализация интерфейса Service.
type service struct {
	teamRepo        repository.TeamRepository
	userRepo        repository.UserRepository
	pullRequestRepo repository.PullRequestRepository

	rndMu sync.Mutex
	rnd   *rand.Rand
}

// NewService создаёт новый экземпляр Service.
func NewService(
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
	pullRequestRepo repository.PullRequestRepository,
) Service {
	return &service{
		teamRepo:        teamRepo,
		userRepo:        userRepo,
		pullRequestRepo: pullRequestRepo,
		rnd:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
