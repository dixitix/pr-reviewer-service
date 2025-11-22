// Package service описывает интерфейсы и доменные ошибки бизнес-логики сервиса назначения ревьюеров.
package service

import "errors"

// Базовые доменные ошибки сервиса.
var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrNotFound          = errors.New("resource not found")
)
