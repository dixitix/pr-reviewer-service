// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import "github.com/dixitix/pr-reviewer-service/internal/domain"

// mapTeamDTOToDomain конвертирует HTTP-DTO команды в доменную команду и её участников.
func mapTeamDTOToDomain(dto TeamDTO) (domain.Team, []domain.User) {
	team := domain.Team{
		Name: domain.TeamName(dto.TeamName),
	}

	members := make([]domain.User, len(dto.Members))
	for i, m := range dto.Members {
		members[i] = domain.User{
			ID:       domain.UserID(m.UserID),
			Username: m.Username,
			TeamName: team.Name,
			IsActive: m.IsActive,
		}
	}

	return team, members
}

// mapTeamDomainToDTO конвертирует доменную команду и её участников в HTTP-DTO.
func mapTeamDomainToDTO(team domain.Team, users []domain.User) TeamDTO {
	members := make([]TeamMemberDTO, len(users))

	for i, u := range users {
		members[i] = TeamMemberDTO{
			UserID:   string(u.ID),
			Username: u.Username,
			IsActive: u.IsActive,
		}
	}

	return TeamDTO{
		TeamName: string(team.Name),
		Members:  members,
	}
}
