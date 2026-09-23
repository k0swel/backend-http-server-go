package utils

import (
	"encoding/json"
	"io"
)

// Основное назначение функции - декодирование байтов в JSON из тела клиентского запроса
func BodyToJSON(Body io.ReadCloser, payload_output any) error {
	err := json.NewDecoder(Body).Decode(payload_output)
	return err
}
