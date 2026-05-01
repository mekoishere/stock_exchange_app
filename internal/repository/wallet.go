package repository

import (
	"context"
	"database/sql"
	"fmt"
	"stock_exchange_app/internal/models"
)

func GetWalletStocks(walletID string) ([]models.WalletStock, error) {
	var stocks []models.WalletStock

	rows, err := DB.Query("SELECT stock_name, quantity FROM wallet_stocks WHERE wallet_id = ?", walletID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

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

func ExecuteTrade(walletID, stockName, tradeType string) error {
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {

		}
	}(tx)

	var bankQuantity int
	err = tx.QueryRow("SELECT quantity FROM bank_stocks WHERE stock_name = ?", stockName).Scan(&bankQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("STOCK_NOT_FOUND")
		}
		return err
	}

	if tradeType == "buy" {
		if bankQuantity <= 0 {
			return fmt.Errorf("INSUFFICIENT_BANK_STOCKS")
		}

		_, err = tx.Exec("UPDATE bank_stocks SET quantity = quantity - 1 WHERE stock_name = ?", stockName)
		if err != nil {
			return err
		}

		_, _ = tx.Exec("INSERT IGNORE INTO wallets (id) VALUES (?)", walletID)
		_, err = tx.Exec(`
			INSERT INTO wallet_stocks (wallet_id, stock_name, quantity) 
			VALUES (?, ?, 1) 
			ON DUPLICATE KEY UPDATE quantity = quantity + 1`,
			walletID, stockName)
		if err != nil {
			return err
		}

	} else if tradeType == "sell" {
		var walletQuantity int
		err := tx.QueryRow("SELECT quantity FROM wallet_stocks WHERE wallet_id = ? AND stock_name = ?", walletID, stockName).Scan(&walletQuantity)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("INSUFFICIENT_WALLET_STOCKS")
			}
			return err
		}

		if walletQuantity <= 0 {
			return fmt.Errorf("INSUFFICIENT_WALLET_STOCKS")
		}

		if walletQuantity == 1 {
			_, err = tx.Exec("DELETE FROM wallet_stocks WHERE wallet_id = ? AND stock_name = ?", walletID, stockName)
		} else {
			_, err = tx.Exec("UPDATE wallet_stocks SET quantity = quantity - 1 WHERE wallet_id = ? AND stock_name = ?", walletID, stockName)
		}
		if err != nil {
			return err
		}

		_, err = tx.Exec("UPDATE bank_stocks SET quantity = quantity + 1 WHERE stock_name = ?", stockName)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("INVALID_TRADE_TYPE")
	}

	if _, err := tx.Exec("INSERT INTO audit_log (operation_type, wallet_id, stock_name) VALUES (?, ?, ?)", tradeType, walletID, stockName); err != nil {
		return err
	}

	return tx.Commit()
}
