# Simplified Stock Exchange Simulation

A high-performance, containerized, and highly available stock market simulation service. This project demonstrates best engineering practices in Go, including transactional database management, clean architecture, and infrastructure redundancy.

## Key Features

- **Transactional Integrity**: Every buy/sell operation is wrapped in a strict SQL transaction. We ensure that a stock is never removed from the bank without being added to a wallet (and vice-versa).
- **High Availability (HA)**: Designed for resilience. The system runs multiple application replicas behind an Nginx Load Balancer. Even if an instance is terminated, the service remains reachable.
- **Automated Infrastructure**: Single-command startup using Docker Compose, including health checks to ensure the Go application only starts when the MySQL database is ready.
- **Audit Logging**: A persistent record of all successful user transactions, ordered chronologically.
- **Zero-Friction Setup**: Cross-platform startup scripts (Windows/Linux/macOS) that handle environment configuration and port mapping.

## Architecture

- **Language**: Go (Golang) 1.22+
- **Database**: MySQL 8.0 (with automated schema initialization)
- **Load Balancer**: Nginx (configured as a reverse proxy with failover logic)
- **Design Pattern**: Clean Architecture (Separation of concerns between API, Repository, and Models)

The traffic flow is as follows:
`User Request` -> `Nginx (Port XXXX)` -> `App Instance 1 or 2` -> `MySQL DB`

##  Getting Started

You only need **Docker** and **Docker Compose** installed. No other dependencies (like Go or MySQL) are required locally as they are handled within containers.

### One-Command Startup

Provide the desired port as an argument to the startup script:

**On Windows (PowerShell/CMD):**
```powershell
.\run.bat 8080
```

**On Linux / macOS:**
```bash
chmod +x run.sh
./run.sh 8080
```
*The service will be available at `http://localhost:8080`.*

## API Reference

### Bank Management
- `GET /stocks`: Returns current bank inventory.
- `POST /stocks`: Resets the bank state and clears all user wallets (Body: `{"stocks": [{"name": "AAPL", "quantity": 100}]}`).

### Trading
- `POST /wallets/{wallet_id}/stocks/{stock_name}`: Buy or sell a single stock.
  - Body: `{"type": "buy"}` or `{"type": "sell"}`.
  - Returns `200 OK` on success, `404` if stock is unknown, or `400` if bank/wallet has insufficient quantity.
  - Automatically creates a wallet if it doesn't exist.

### Wallet Information
- `GET /wallets/{wallet_id}`: Returns the current state of a specific wallet.
- `GET /wallets/{wallet_id}/stocks/{stock_name}`: Returns the quantity of a specific stock in a wallet.

### System & Audit
- `GET /log`: Returns the complete audit log of all successful trade operations.
- `GET /health`: Returns the health status of the current instance.
- `POST /chaos`: Simulates a failure by terminating the current instance (useful for testing HA).

## Testing High Availability

1. Start the project on port `8080`.
2. Open a browser or tool like Bruno/Postman and call `GET http://localhost:8080/health` (Status: 200).
3. Call `POST http://localhost:8080/chaos`. You will see `appX exited with code 1` in your terminal.
4. Immediately call `GET http://localhost:8080/health` again. 
5. **Result**: You will still get a `200 OK` response. Nginx detected the failure of the first instance and seamlessly rerouted your request to the second healthy instance.

## Project Structure

```text
├── cmd/server/       # Application entry point
├── internal/
│   ├── api/          # HTTP Handlers and request parsing
│   ├── repository/   # Database logic and transactions
│   └── models/       # Shared data structures
├── Dockerfile        # Multi-stage build for Go
├── compose.yaml      # Orchestration for HA (2x App, Nginx, MySQL)
├── nginx.conf        # Load balancer and failover configuration
└── create-tables.sql # Database schema
```
