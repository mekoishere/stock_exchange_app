package main

import (
	"log"
	"net/http"
	"os"
	"stock_exchange_app/internal/api"
	"stock_exchange_app/internal/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if err := repository.InitDB(); err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/wallets/{wallet_id}", api.GetWalletByIdHandler).Methods("GET")
	r.HandleFunc("/wallets/{wallet_id}/stocks/{stock_name}", api.PostTradeHandler).Methods("POST")
	r.HandleFunc("/wallets/{wallet_id}/stocks/{stock_name}", api.GetStockInWalletHandler).Methods("GET")
	r.HandleFunc("/stocks", api.GetBankStocksHandler).Methods("GET")
	r.HandleFunc("/stocks", api.PostBankStocksHandler).Methods("POST")
	r.HandleFunc("/log", api.GetAuditLogsHandler).Methods("GET")
	r.HandleFunc("/chaos", api.ChaosHandler).Methods("POST")

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "2137"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
