// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

// TeamMemberDTO представляет участника команды в HTTP-слое.
type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// TeamDTO представляет команду и её участников в HTTP-слое.
type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

// GetTeamResponse описывает ответ на запрос получения команды.
type GetTeamResponse struct {
	Team TeamDTO `json:"team"`
}
