-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Users Table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Stocks Table
CREATE TABLE stocks (
    symbol VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    base_currency VARCHAR(10) DEFAULT 'INR'
);

-- 3. Stock Prices Table (Historical)
CREATE TABLE stock_prices (
    id BIGSERIAL PRIMARY KEY,
    stock_symbol VARCHAR(50) REFERENCES stocks(symbol),
    price NUMERIC(18, 4) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_stock_prices_symbol_time ON stock_prices(stock_symbol, timestamp DESC);

-- 4. Rewards Table
CREATE TABLE rewards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id BIGINT REFERENCES users(id),
    stock_symbol VARCHAR(50) REFERENCES stocks(symbol),
    quantity NUMERIC(18, 6) NOT NULL CHECK (quantity > 0),
    awarded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reference_id VARCHAR(255) UNIQUE NOT NULL  -- For Idempotency
);

CREATE INDEX idx_rewards_user_id ON rewards(user_id);
CREATE INDEX idx_rewards_awarded_at ON rewards(awarded_at);

-- 5. Ledger Accounts Table
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(20) NOT NULL -- ASSET, LIABILITY, EXPENSE, INCOME, EQUITY
);

-- Seed basic accounts
INSERT INTO accounts (name, type) VALUES 
('Bank', 'ASSET'),
('Reward_Expense', 'EXPENSE'),
('Brokerage_Fees', 'EXPENSE'),
('Company_Stock_Inventory', 'ASSET'),
('User_Stock_Liability', 'LIABILITY'); 
-- Note: 'User_Stock_Liability' represents the obligation to the user (the shares they own)

-- 6. Journal Entries (Transaction Header)
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference_type VARCHAR(50) NOT NULL, -- e.g. 'REWARD'
    reference_id VARCHAR(255), -- ID of the reward or external event
    description TEXT,
    posted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 7. Journal Postings (Transaction Lines)
CREATE TABLE journal_postings (
    id BIGSERIAL PRIMARY KEY,
    journal_entry_id UUID REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id INT REFERENCES accounts(id),
    amount NUMERIC(18, 6) NOT NULL CHECK (amount >= 0),
    direction VARCHAR(10) NOT NULL CHECK (direction IN ('DEBIT', 'CREDIT')),
    asset_type VARCHAR(10) NOT NULL CHECK (asset_type IN ('INR', 'STOCK')),
    stock_symbol VARCHAR(50) REFERENCES stocks(symbol), -- Nullable if INR
    
    CONSTRAINT check_stock_symbol CHECK (
        (asset_type = 'STOCK' AND stock_symbol IS NOT NULL) OR 
        (asset_type = 'INR' AND stock_symbol IS NULL)
    )
);

CREATE INDEX idx_postings_entry_id ON journal_postings(journal_entry_id);
CREATE INDEX idx_postings_account_id ON journal_postings(account_id);
