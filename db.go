package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB *sqlx.DB

func ConnectDB(dsn string) {
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	if err := DB.Ping(); err != nil {
		logrus.Fatalf("Failed to ping database: %v", err)
	}
	logrus.Info("Connected to database")
}

func CreateReward(req CreateRewardRequest, price float64) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Idempotency
	var exists int
	err = tx.Get(&exists, "SELECT count(*) FROM rewards WHERE reference_id = $1", req.ReferenceID)
	if err != nil {
		return err
	}
	if exists > 0 {
		return fmt.Errorf("duplicate reward reference")
	}

	// User Check: Auto-create if not exists (for testing flexibility)
	var userExists int
	err = tx.Get(&userExists, "SELECT count(*) FROM users WHERE id = $1", req.UserID)
	if err != nil {
		return err
	}
	if userExists == 0 {
		_, err = tx.Exec(`
			INSERT INTO users (id, name, email) 
			VALUES ($1, $2, $3)
		`, req.UserID, fmt.Sprintf("User %d", req.UserID), fmt.Sprintf("user%d@stocky.com", req.UserID))
		if err != nil {
			return fmt.Errorf("failed to auto-create user: %v", err)
		}
	}

	rewardID := uuid.New()
	_, err = tx.Exec(`
		INSERT INTO rewards (id, user_id, stock_symbol, quantity, reference_id, awarded_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, rewardID, req.UserID, req.StockSymbol, req.Quantity, req.ReferenceID, time.Now())
	if err != nil {
		return err
	}

	totalValue := price * req.Quantity
	brokerage := totalValue * 0.01
	totalCashOut := totalValue + brokerage

	entryID := uuid.New()
	_, err = tx.Exec(`
		INSERT INTO journal_entries (id, reference_type, reference_id, description, posted_at)
		VALUES ($1, 'REWARD', $2, $3, $4)
	`, entryID, req.ReferenceID, fmt.Sprintf("Reward %f %s to User %d", req.Quantity, req.StockSymbol, req.UserID), time.Now())
	if err != nil {
		return err
	}

	postings := []JournalPosting{
		{JournalEntryID: entryID, AccountID: 2, Amount: totalValue, Direction: "DEBIT", AssetType: "INR"},
		{JournalEntryID: entryID, AccountID: 3, Amount: brokerage, Direction: "DEBIT", AssetType: "INR"},
		{JournalEntryID: entryID, AccountID: 1, Amount: totalCashOut, Direction: "CREDIT", AssetType: "INR"},
		{JournalEntryID: entryID, AccountID: 5, Amount: req.Quantity, Direction: "CREDIT", AssetType: "STOCK", StockSymbol: &req.StockSymbol},
		{JournalEntryID: entryID, AccountID: 4, Amount: req.Quantity, Direction: "DEBIT", AssetType: "STOCK", StockSymbol: &req.StockSymbol},
	}

	for _, p := range postings {
		_, err = tx.Exec(`
			INSERT INTO journal_postings (journal_entry_id, account_id, amount, direction, asset_type, stock_symbol)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, p.JournalEntryID, p.AccountID, p.Amount, p.Direction, p.AssetType, p.StockSymbol)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetTodayRewards(userID int64) ([]Reward, error) {
	var rewards []Reward
	err := DB.Select(&rewards, "SELECT * FROM rewards WHERE user_id = $1 AND awarded_at >= current_date", userID)
	return rewards, err
}

func GetUserStats(userID int64) (*UserStatsResponse, error) {
	var stockStats []StockStat
	err := DB.Select(&stockStats, `
		SELECT stock_symbol, CAST(SUM(quantity) AS DOUBLE PRECISION) as total_quantity 
		FROM rewards 
		WHERE user_id = $1 AND awarded_at >= current_date
		GROUP BY stock_symbol
	`, userID)
	if err != nil {
		return nil, err
	}

	query := `
		WITH UserHoldings AS (
			SELECT stock_symbol, SUM(quantity) as qty
			FROM rewards
			WHERE user_id = $1
			GROUP BY stock_symbol
		),
		LatestPrices AS (
			SELECT DISTINCT ON (stock_symbol) stock_symbol, price
			FROM stock_prices
			ORDER BY stock_symbol, timestamp DESC
		)
		SELECT COALESCE(CAST(SUM(u.qty * lp.price) AS DOUBLE PRECISION), 0.0)
		FROM UserHoldings u
		JOIN LatestPrices lp ON u.stock_symbol = lp.stock_symbol
	`
	var portfolioValue float64
	err = DB.Get(&portfolioValue, query, userID)
	if err != nil {
		return nil, err
	}

	return &UserStatsResponse{
		TotalSharesToday: stockStats,
		PortfolioValue:   portfolioValue,
	}, nil
}

func GetHistoricalPortfolio(userID int64) ([]map[string]interface{}, error) {
	query := `
		WITH dates AS (
			SELECT generate_series(current_date - interval '30 days', current_date, '1 day')::date AS day
		),
		daily_holdings AS (
			SELECT d.day, r.stock_symbol, SUM(r.quantity) as cum_qty
			FROM dates d
			JOIN rewards r ON r.awarded_at::date <= d.day
			WHERE r.user_id = $1
			GROUP BY d.day, r.stock_symbol
		),
		daily_prices AS (
			SELECT 
				sp.stock_symbol, 
				sp.timestamp::date as day, 
				sp.price,
				ROW_NUMBER() OVER(PARTITION BY sp.stock_symbol, sp.timestamp::date ORDER BY sp.timestamp DESC) as rn
			FROM stock_prices sp
		)
		SELECT dh.day, CAST(SUM(dh.cum_qty * dp.price) AS DOUBLE PRECISION) as value
		FROM daily_holdings dh
		JOIN daily_prices dp ON dh.stock_symbol = dp.stock_symbol AND dh.day = dp.day AND dp.rn = 1
		GROUP BY dh.day
		ORDER BY dh.day
	`
	type HistoryPoint struct {
		Day   string  `db:"day"`
		Value float64 `db:"value"`
	}
	// Initialize as empty slice to avoid returning null in JSON
	history := make([]HistoryPoint, 0)
	err := DB.Select(&history, query, userID)

	var result []map[string]interface{}
	// Manual mapping isn't strictly necessary if HistoryPoint had json tags,
	// but keeping existing logic structure
	if len(history) == 0 {
		return []map[string]interface{}{}, nil
	}
	for _, h := range history {
		result = append(result, map[string]interface{}{
			"date":  h.Day,
			"value": h.Value,
		})
	}
	return result, err
}

func GetPortfolio(userID int64) ([]PortfolioItem, error) {
	query := `
		WITH UserHoldings AS (
			SELECT stock_symbol, SUM(quantity) as qty
			FROM rewards
			WHERE user_id = $1
			GROUP BY stock_symbol
		),
		LatestPrices AS (
			SELECT DISTINCT ON (stock_symbol) stock_symbol, price
			FROM stock_prices
			ORDER BY stock_symbol, timestamp DESC
		)
		SELECT u.stock_symbol, CAST(u.qty AS DOUBLE PRECISION) as total_quantity, 
		       CAST(lp.price AS DOUBLE PRECISION) as current_price, 
		       CAST(u.qty * lp.price AS DOUBLE PRECISION) as value_inr
		FROM UserHoldings u
		JOIN LatestPrices lp ON u.stock_symbol = lp.stock_symbol
	`
	var items []PortfolioItem
	err := DB.Select(&items, query, userID)
	return items, err
}

func InsertStockPrice(symbol string, price float64) error {
	_, err := DB.Exec(`INSERT INTO stock_prices (stock_symbol, price, timestamp) VALUES ($1, $2, $3)`, symbol, price, time.Now())
	return err
}
