package models

type WalletStock struct {
	StockName string `json:"name"`
	Quantity  int    `json:"quantity"`
}

type WalletResponse struct {
	ID     string        `json:"id"`
	Stocks []WalletStock `json:"stocks"`
}

type BankStock struct {
	StockName string `json:"name"`
	Quantity  int    `json:"quantity"`
}

type AuditLogEntry struct {
	Type      string `json:"type"`
	WalletID  string `json:"wallet_id"`
	StockName string `json:"stock_name"`
}
