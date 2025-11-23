// Package user содержит обработчики и DTO для работы с пользователями.
package user

import "github.com/dixitix/pr-reviewer-service/internal/domain"

// mapUserDomainToDTO конвертирует доменного пользователя в HTTP-DTO.
func mapUserDomainToDTO(u domain.User) DTO {
	return DTO{
		UserID:   string(u.ID),
		Username: u.Username,
		TeamName: string(u.TeamName),
		IsActive: u.IsActive,
	}
}
