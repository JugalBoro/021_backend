# Stocky Backend

Backend service for Stocky, a platform rewarding users with stock shares.

## Features

- **Reward System**: Assign fractional shares to users.
- **Double-Entry Ledger**: Tracks every financial and stock movement (Assets, Liabilities, Expenses).
- **Portfolio Tracking**: Real-time valuation of user holdings.
- **Historical Data**: Track portfolio value over time.
- **Background Jobs**: Updates stock prices hourly.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **Libraries**:
  - `sqlx`: Database extensions
  - `logrus`: Structured logging
  - `uuid`: UUID generation

## Setup & Running

### Prerequisites

- Go 1.20+
- PostgreSQL

### Database Setup

1. Create a database named `assignment`.
2. Run the SQL script in `db/schema.sql` to initialize tables and seed data.

```bash
psql -U postgres -d assignment -f db/schema.sql
```

### Running the Application

1. Clone the repository.
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set environment variables (optional, defaults are set in `main.go`):
   - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
4. Run the server:
   ```bash
   go run .
   ```

## Project Structure

The project has been simplified for ease of understanding and development.

- `main.go`: Application entry point, config, and router.
- `handlers.go`: API endpoint logic.
- `db.go`: Database connection and queries.
- `models.go`: Data structures (Entities & DTOs).
- `cron.go`: Background jobs (Price Updater).
- `migrations/`: SQL schema files.


