package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stock_exchange_app/internal/models"
	"stock_exchange_app/internal/repository"

	"github.com/gorilla/mux"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	if err != nil {
		return
	}
}

func GetWalletByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["wallet_id"]

	stocks, err := repository.GetWalletStocks(walletID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := models.WalletResponse{
		ID:     walletID,
		Stocks: stocks,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func GetStockInWalletHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["wallet_id"]
	stockName := vars["stock_name"]

	quantity, err := repository.GetStockQuantityInWallet(walletID, stockName)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%d", quantity)
}

func PostTradeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID := vars["wallet_id"]
	stockName := vars["stock_name"]

	var req models.TradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Type != "buy" && req.Type != "sell" {
		http.Error(w, "Type must be 'buy' or 'sell'", http.StatusBadRequest)
		return
	}

	err := repository.ExecuteTrade(walletID, stockName, req.Type)
	if err != nil {
		switch err.Error() {
		case "STOCK_NOT_FOUND":
			http.Error(w, "Stock not found", http.StatusNotFound)
			return

		case "INSUFFICIENT_BANK_STOCKS":
			http.Error(w, "Insufficient funds", http.StatusBadRequest)
			return

		case "INSUFFICIENT_WALLET_STOCKS":
			http.Error(w, "Insufficient wallet stocks", http.StatusBadRequest)
			return

		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
