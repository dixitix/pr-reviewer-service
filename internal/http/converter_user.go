package httpserver

import "github.com/dixitix/pr-reviewer-service/internal/domain"

// mapUserDomainToDTO конвертирует доменного пользователя в HTTP-DTO.
func mapUserDomainToDTO(u domain.User) UserDTO {
	return UserDTO{
		UserID:   string(u.ID),
		Username: u.Username,
		TeamName: string(u.TeamName),
		IsActive: u.IsActive,
	}
}
