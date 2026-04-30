-- Drop existing database if it exists and create the new one
DROP DATABASE IF EXISTS stock_market;
CREATE DATABASE stock_market;
USE stock_market;

-- 1. User Wallets
-- Tracks unique wallet IDs.
CREATE TABLE wallets (
                         id VARCHAR(255) PRIMARY KEY
);

-- 2. User Holdings
-- Tracks quantity of specific stocks owned by each wallet.
CREATE TABLE wallet_stocks (
                               wallet_id VARCHAR(255) NOT NULL,
                               stock_name VARCHAR(255) NOT NULL,
                               quantity INT NOT NULL DEFAULT 0,
                               PRIMARY KEY (wallet_id, stock_name),
                               FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);

-- 3. Bank Inventory
-- Dedicated table for the Bank's stock availability.
CREATE TABLE bank_stocks (
                             stock_name VARCHAR(255) PRIMARY KEY,
                             quantity INT NOT NULL DEFAULT 0
);

-- 4. Audit Log
-- Records all successful buy/sell operations on user wallets.
CREATE TABLE audit_log (
                           id INT AUTO_INCREMENT PRIMARY KEY,
                           operation_type ENUM('buy', 'sell') NOT NULL,
                           wallet_id VARCHAR(255) NOT NULL,
                           stock_name VARCHAR(255) NOT NULL,
                           created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6)
);

-- Create a user for the application (optional, aligned with previous script)
-- Note: 'username' and 'password' should be changed for production use.
CREATE USER IF NOT EXISTS 'username' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON stock_market.* TO 'username'@'%';
FLUSH PRIVILEGES;
