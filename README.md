# Stocky Backend

Backend service for Stocky, a platform that rewards users with Indian stock shares. This system handles reward processing, double-entry ledger accounting, stock price tracking, and portfolio valuation.

## Features

- **Reward Processing**: Idempotent API to issue stock rewards.
- **Double-Entry Ledger**: Tracks all transactions (INR & Stock) with proper accounting principles (Assets, Liabilities, Expenses).
- **Real-time Price Updates**: Background worker fetches stock prices every hour.
- **Portfolio Analytics**: APIs for daily stats, historical valuation charts, and current portfolio breakdown.
- **Simplified Architecture**: Easy to navigate flat file structure.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin (HTTP Web Framework)
- **Database**: PostgreSQL
- **ORM/SQL**: SQLx (Extensions to database/sql)
- **Logging**: Logrus

## Prerequisites

- Go 1.20+
- PostgreSQL
- Git

## Setup & Running

1. **Clone the repository**:
   ```bash
   git clone https://github.com/JugalBoro/021_backend.git
   cd 021_backend
   ```

2. **Database Setup**:
   - Create a PostgreSQL database named `assignment`.
   - Run the schema migration:
     ```bash
     psql -U postgres -d assignment -f migrations/schema.sql
     ```
   *(Note: The schema script automatically seeds necessary data like supported stocks and account types.)*

3. **Environment Configuration**:
   - Create a `.env` file in the root directory:
     ```env
     DB_HOST=localhost
     DB_PORT=5432
     DB_USER=postgres
     DB_PASSWORD=postgres
     DB_NAME=assignment
     PORT=8080
     ```

4. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

5. **Run the Server**:
   ```bash
   go run .
   ```

## API Endpoints

### 1. Create Reward
- **POST** `/api/reward`
- **Body**:
  ```json
  {
      "user_id": 1,
      "stock_symbol": "RELIANCE",
      "quantity": 1.5,
      "reference_id": "unique-ref-123"
  }
  ```

### 2. Get Today's Rewards
- **GET** `/api/today-stocks/:userId`

### 3. Get User Stats
- **GET** `/api/stats/:userId`
- Returns total shares received today and total portfolio value.

### 4. Get Historical Portfolio Value
- **GET** `/api/historical-inr/:userId`
- Returns a list of `{date, value}` for the last 30 days.

### 5. Get Portfolio Breakdown
- **GET** `/api/portfolio/:userId`
- Returns holdings per stock with current market price and value.

## Project Structure

- `main.go`: Application entry point, config loading, router setup.
- `handlers.go`: API request handlers.
- `db.go`: Database connection and data access logic.
- `models.go`: Go structs for Database entities and API DTOs.
- `cron.go`: Background job for updating stock prices.
- `migrations/`: SQL schema and seed data.
- `docs/`: Deployment docs, API specs, and Postman collection.

## Deployment

This application can be deployed on any platform supporting Go binaries (e.g., Render, Heroku, AWS EC2, or Docker).
