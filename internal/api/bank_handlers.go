package api

import (
	"encoding/json"
	"net/http"
	"stock_exchange_app/internal/models"
	"stock_exchange_app/internal/repository"
)

func GetBankStocksHandler(w http.ResponseWriter, r *http.Request) {
	stocks, err := repository.GetBankStocks()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if stocks == nil {
		stocks = []models.BankStock{}
	}

	response := map[string]interface{}{
		"stocks": stocks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PostBankStocksHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Stocks []models.BankStock `json:"stocks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := repository.UpdateBankStocks(req.Stocks); err != nil {
		http.Error(w, "Failed to update bank stocks", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
