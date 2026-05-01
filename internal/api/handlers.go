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
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
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
	json.NewEncoder(w).Encode(response)
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
