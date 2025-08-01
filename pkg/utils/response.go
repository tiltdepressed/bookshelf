// Package utils
package utils

import (
	"encoding/json"
	"net/http"
)

// Вспомогательная функция для обработки JSON

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
