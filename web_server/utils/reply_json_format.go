package utils

import (
	"encoding/json"
	"net/http"
	"web_server_go/utils"
)

// Основное назначение функции - отправка ответа клиенту в JSON (валидный ответ без ошибок)
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		utils.Logger.Printf("writeJSON: %v", err)
	}
}

// Основное назначение функции - отправка ответа клиенту в JSON (В случае какой-либо ошибки)
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}
