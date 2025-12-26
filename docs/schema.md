# Database Schema Design

## Overview
The database is designed to handle user rewards, stock data, and a double-entry ledger system for financial and inventory tracking.
The database name will be `assignment`.

## Tables

### 1. `users`
Stores user information.
- `id` (BIGSERIAL, PK)
- `name` (VARCHAR)
- `created_at` (TIMESTAMP)

### 2. `stocks`
Stores supported stock symbols.
- `symbol` (VARCHAR, PK) - e.g., 'RELIANCE', 'TCS'
- `name` (VARCHAR)
- `base_currency` (VARCHAR) - Default 'INR'

### 3. `stock_prices`
Stores historical stock prices for valuation.
- `id` (BIGSERIAL, PK)
- `stock_symbol` (VARCHAR, FK -> stocks.symbol)
- `price` (NUMERIC(18, 4))
- `timestamp` (TIMESTAMP)
- Indexes: `(stock_symbol, timestamp)`

### 4. `rewards`
Records the event of a user being rewarded. This is the source of truth for "User Portfolio".
- `id` (UUID, PK) - To prevent enumeration and aid idempotency
- `user_id` (BIGINT, FK -> users.id)
- `stock_symbol` (VARCHAR, FK -> stocks.symbol)
- `quantity` (NUMERIC(18, 6)) - Allows fractional shares
- `awarded_at` (TIMESTAMP)
- `reference_id` (VARCHAR, UNIQUE) - External ID to prevent replay attacks (idempotency key)

### 5. `accounts`
Defines the accounts for the double-entry ledger.
- `id` (SERIAL, PK)
- `name` (VARCHAR) - e.g., 'Bank', 'Reward_Expense', 'Brokerage_Fees', 'User_Liability'
- `type` (VARCHAR) - ASSET, LIABILITY, EQUITY, INCOME, EXPENSE

### 6. `journal_entries`
Groups ledger postings into a single atomic transaction.
- `id` (UUID, PK)
- `reference_type` (VARCHAR) - e.g., 'REWARD_ISSUANCE'
- `reference_id` (VARCHAR) - Links to rewards.id or other event IDs
- `description` (TEXT)
- `posted_at` (TIMESTAMP)

### 7. `journal_postings`
Individual debit/credit lines.
- `id` (BIGSERIAL, PK)
- `journal_entry_id` (UUID, FK -> journal_entries.id)
- `account_id` (INT, FK -> accounts.id)
- `amount` (NUMERIC(18, 6)) - Always positive
- `direction` (VARCHAR) - 'DEBIT' or 'CREDIT'
- `asset_type` (VARCHAR) - 'INR' or 'STOCK'
- `stock_symbol` (VARCHAR, Nullable, FK -> stocks.symbol) - Populated if asset_type is STOCK

## Relationships & Logic

- **Reward Event**:
  When a user gets a reward, a row is inserted into `rewards`.
  Simultaneously, a `journal_entry` is created with multiple `journal_postings`:
  1. **Financial Leg** (INR):
     - Dr `Reward_Expense` (Price * Qty)
     - Dr `Brokerage_Fees` (Fee)
     - Cr `Bank` (Total Outflow)
  2. **Inventory Leg** (Stock Units):
     - Dr `User_Holdings_Contra` (Liability/Tracking) - Qty
     - Cr `Company_Inventory` - Qty

## Data Types
- Stock Quantities: `NUMERIC(18, 6)`
- INR/Price: `NUMERIC(18, 4)`

## Indexes
- `rewards`: `user_id`, `awarded_at` (for fast history lookups)
- `stock_prices`: `stock_symbol`, `timestamp` DESC (for getting latest price)
