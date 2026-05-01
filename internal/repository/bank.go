package repository

import (
	"context"
	"stock_exchange_app/internal/models"
)

func GetBankStocks() ([]models.BankStock, error) {
	var stocks []models.BankStock

	rows, err := DB.Query("SELECT stock_name, quantity FROM bank_stocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.BankStock
		if err := rows.Scan(&s.StockName, &s.Quantity); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	return stocks, nil
}

func UpdateBankStocks(stocks []models.BankStock) error {
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM bank_stocks"); err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM wallet_stocks"); err != nil {
		return err
	}

	for _, s := range stocks {
		if _, err := tx.Exec("INSERT INTO bank_stocks (stock_name, quantity) VALUES (?, ?)", s.StockName, s.Quantity); err != nil {
			return err
		}
	}

	return tx.Commit()
}
