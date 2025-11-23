// Package httperr содержит структуры и функции формирования HTTP-ошибок.
package httperr

import (
	"encoding/json"
	"log"
	"net/http"
)

// HTTP-коды ошибок.
const (
	ErrorCodeTeamExists  = "TEAM_EXISTS"
	ErrorCodePRExists    = "PR_EXISTS"
	ErrorCodePRMerged    = "PR_MERGED"
	ErrorCodeNotAssigned = "NOT_ASSIGNED"
	ErrorCodeNoCandidate = "NO_CANDIDATE"
	ErrorCodeNotFound    = "NOT_FOUND"
)

// ErrorResponseBody описывает тело ошибки в формате ErrorResponse.
type ErrorResponseBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse описывает стандартный ответ об ошибке сервиса.
type ErrorResponse struct {
	Error ErrorResponseBody `json:"error"`
}

// WriteJSONError отправляет ответ об ошибке в формате ErrorResponse.
// Ошибка записи логируется через переданный логгер, при его отсутствии используется стандартный лог.
func WriteJSONError(w http.ResponseWriter, statusCode int, code, message string, logger *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Error: ErrorResponseBody{
			Code:    code,
			Message: message,
		},
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		if logger != nil {
			logger.Printf("failed to write error response: %v", err)
			return
		}

		log.Printf("failed to write error response: %v", err)
	}
}
