package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type WalletStock struct {
	WalletId  string `json:"-"`
	StockName string `json:"name"`
	Quantity  int    `json:"quantity"`
}

type WalletResponse struct {
	ID     string        `json:"id"`
	Stocks []WalletStock `json:"stocks"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func GetWalletByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var walletStocks []WalletStock

	query, quErr := db.Query("SELECT * FROM wallet_stocks WHERE wallet_id = ?", vars["wallet_id"])
	if quErr != nil {
		log.Println(quErr)
		return
	}
	defer query.Close()

	for query.Next() {
		var walletStock WalletStock
		if err := query.Scan(&walletStock.WalletId, &walletStock.StockName, &walletStock.Quantity); err != nil {
			log.Println(err)
			return
		}
		walletStocks = append(walletStocks, walletStock)
	}

	response := WalletResponse{
		ID:     vars["wallet_id"],
		Stocks: walletStocks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

var db *sql.DB

func main() {
	dotEnvErr := godotenv.Load()
	if dotEnvErr != nil {
		log.Fatal("Error loading .env file")
	}
	// Capture connection properties.
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:" + os.Getenv("DBPORT")
	cfg.DBName = "stock_market"

	var sqlOpenErr error
	db, sqlOpenErr = sql.Open("mysql", cfg.FormatDSN())
	if sqlOpenErr != nil {
		log.Fatal(sqlOpenErr)
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheckHandler)
	r.HandleFunc("/wallets/{wallet_id}", GetWalletByIdHandler)

	serverListenErr := http.ListenAndServe(":2137", r)
	if serverListenErr != nil {
		return
	}
}
