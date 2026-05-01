package repository

import (
	"database/sql"
	"stock_exchange_app/internal/models"
)

func GetWalletStocks(walletID string) ([]models.WalletStock, error) {
	var stocks []models.WalletStock

	rows, err := DB.Query("SELECT stock_name, quantity FROM wallet_stocks WHERE wallet_id = ?", walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.WalletStock
		if err := rows.Scan(&s.StockName, &s.Quantity); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	return stocks, nil
}

func GetStockQuantityInWallet(walletID, stockName string) (int, error) {
	var quantity int
	err := DB.QueryRow("SELECT quantity FROM wallet_stocks WHERE wallet_id = ? AND stock_name = ?", walletID, stockName).Scan(&quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // If no row, user has 0 stocks
		}
		return 0, err
	}
	return quantity, nil
}
