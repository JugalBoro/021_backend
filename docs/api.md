# API Specification

## 1. POST /reward
Records that a user has been rewarded shares.

**Endpoint:** `POST /reward`

**Request Body:**
```json
{
  "user_id": 123,
  "stock_symbol": "RELIANCE",
  "quantity": "1.500000",
  "timestamp": "2024-12-26T10:00:00Z",
  "reference_id": "evt_001"
}
```

**Response (201 Created):**
```json
{
  "message": "Reward processed successfully",
  "reward_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
}
```

**Edge Cases:**
- Returns `409 Conflict` if `reference_id` already exists (Idempotency).
- Returns `400 Bad Request` if stock symbol is invalid.

## 2. GET /today-stocks/{userId}
Return all stock rewards for the user for today (00:00 to 23:59 server time).

**Endpoint:** `GET /today-stocks/:userId`

**Response (200 OK):**
```json
{
  "user_id": 123,
  "date": "2024-12-26",
  "rewards": [
    {
      "stock_symbol": "RELIANCE",
      "quantity": "1.500000",
      "awarded_at": "2024-12-26T10:00:00Z"
    },
    {
      "stock_symbol": "TCS",
      "quantity": "0.500000",
      "awarded_at": "2024-12-26T12:30:00Z"
    }
  ]
}
```

## 3. GET /historical-inr/{userId}
Return the INR value of the user’s stock rewards for all past days (up to yesterday).
Calculated using the closing price of each day.

**Endpoint:** `GET /historical-inr/:userId`

**Response (200 OK):**
```json
{
  "user_id": 123,
  "history": [
    {
      "date": "2024-12-24",
      "portfolio_value_inr": "4500.0000"
    },
    {
      "date": "2024-12-25",
      "portfolio_value_inr": "4650.5000"
    }
  ]
}
```

## 4. GET /stats/{userId}
Returns the stats for today and current portfolio value.

**Endpoint:** `GET /stats/:userId`

**Response (200 OK):**
```json
{
  "user_id": 123,
  "today_rewards_total": {
    "RELIANCE": "1.500000",
    "TCS": "0.500000"
  },
  "current_portfolio_value_inr": "12500.7500"
}
```

## Bonus: GET /portfolio/{userId}
Current holdings per stock.

**Endpoint:** `GET /portfolio/:userId`

**Response:**
```json
{
  "holdings": [
    {
      "stock_symbol": "RELIANCE",
      "total_quantity": "10.000000",
      "current_price": "2500.00",
      "current_value": "25000.0000"
    }
  ]
}
```
