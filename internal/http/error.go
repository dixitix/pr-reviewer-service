// Package httpserver содержит структуры данных и вспомогательные типы
// для HTTP-слоя сервиса назначения ревьюеров.
package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
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

// writeJSONError отправляет ответ об ошибке в формате ErrorResponse.
func writeJSONError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Error: ErrorResponseBody{
			Code:    code,
			Message: message,
		},
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write error response: %v", err)
	}
}
