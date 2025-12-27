package main

import (
	"time"

	"github.com/google/uuid"
)

// Entities representing database tables

type User struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Stock struct {
	Symbol       string `db:"symbol" json:"symbol"`
	Name         string `db:"name" json:"name"`
	BaseCurrency string `db:"base_currency" json:"base_currency"`
}

type StockPrice struct {
	ID          int64     `db:"id" json:"id"`
	StockSymbol string    `db:"stock_symbol" json:"stock_symbol"`
	Price       float64   `db:"price" json:"price"`
	Timestamp   time.Time `db:"timestamp" json:"timestamp"`
}

type Reward struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	StockSymbol string    `db:"stock_symbol" json:"stock_symbol"`
	Quantity    float64   `db:"quantity" json:"quantity"`
	AwardedAt   time.Time `db:"awarded_at" json:"awarded_at"`
	ReferenceID string    `db:"reference_id" json:"reference_id"`
}

type Account struct {
	ID   int32  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Type string `db:"type" json:"type"`
}

type JournalEntry struct {
	ID            uuid.UUID `db:"id" json:"id"`
	ReferenceType string    `db:"reference_type" json:"reference_type"`
	ReferenceID   string    `db:"reference_id" json:"reference_id"`
	Description   string    `db:"description" json:"description"`
	PostedAt      time.Time `db:"posted_at" json:"posted_at"`
}

type JournalPosting struct {
	ID             int64     `db:"id" json:"id"`
	JournalEntryID uuid.UUID `db:"journal_entry_id" json:"journal_entry_id"`
	AccountID      int32     `db:"account_id" json:"account_id"`
	Amount         float64   `db:"amount" json:"amount"`
	Direction      string    `db:"direction" json:"direction"`
	AssetType      string    `db:"asset_type" json:"asset_type"`
	StockSymbol    *string   `db:"stock_symbol" json:"stock_symbol,omitempty"`
}

// Data Transfer Objects (DTOs) for API

type CreateRewardRequest struct {
	UserID      int64   `json:"user_id" binding:"required"`
	StockSymbol string  `json:"stock_symbol" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	ReferenceID string  `json:"reference_id" binding:"required"`
}

type StockStat struct {
	StockSymbol string  `json:"stock_symbol" db:"stock_symbol"`
	TotalQty    float64 `json:"total_quantity" db:"total_quantity"`
}

type UserStatsResponse struct {
	TotalSharesToday []StockStat `json:"total_shares_today"`
	PortfolioValue   float64     `json:"portfolio_value"`
}

type PortfolioItem struct {
	StockSymbol  string  `json:"stock_symbol" db:"stock_symbol"`
	TotalQty     float64 `json:"total_quantity" db:"total_quantity"`
	CurrentPrice float64 `json:"current_price" db:"current_price"`
	ValueINR     float64 `json:"value_inr" db:"value_inr"`
}
