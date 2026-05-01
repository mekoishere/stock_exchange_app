package api

import (
	"encoding/json"
	"net/http"
	"stock_exchange_app/internal/models"
	"stock_exchange_app/internal/repository"
)

func GetAuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := repository.GetAuditLogs()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if logs == nil {
		logs = []models.AuditLogEntry{}
	}

	response := map[string]interface{}{
		"log": logs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
