// Package team содержит обработчики и DTO для работы с командами.
package team

// MemberDTO представляет участника команды в HTTP-слое.
type MemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// DTO представляет команду и её участников в HTTP-слое.
type DTO struct {
	TeamName string      `json:"team_name"`
	Members  []MemberDTO `json:"members"`
}

// GetTeamResponse описывает ответ на запрос получения команды.
type GetTeamResponse struct {
	Team DTO `json:"team"`
}
