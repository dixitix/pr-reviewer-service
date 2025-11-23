// Package user содержит обработчики и DTO для работы с пользователями.
package user

// SetUserActiveRequest описывает тело запроса /users/setIsActive.
type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
