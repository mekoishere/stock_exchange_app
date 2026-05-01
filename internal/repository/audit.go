package repository

import "stock_exchange_app/internal/models"

func GetAuditLogs() ([]models.AuditLogEntry, error) {
	var logs []models.AuditLogEntry

	rows, err := DB.Query("SELECT operation_type, wallet_id, stock_name FROM audit_log ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.AuditLogEntry
		if err := rows.Scan(&entry.Type, &entry.WalletID, &entry.StockName); err != nil {
			return nil, err
		}
		logs = append(logs, entry)
	}
	return logs, nil
}
