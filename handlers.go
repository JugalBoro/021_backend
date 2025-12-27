package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateRewardHandler(c *gin.Context) {
	var req CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch current price
	price, err := GetCurrentPrice(req.StockSymbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock price"})
		return
	}

	err = CreateReward(req, price)
	if err != nil {
		if err.Error() == "duplicate reward reference" {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate reward event"})
			return
		}
		logrus.Error("Failed to create reward: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Reward created successfully"})
}

func GetTodayStocksHandler(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	rewards, err := GetTodayRewards(userID)
	if err != nil {
		logrus.Errorf("GetTodayStocks failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rewards"})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func GetStatsHandler(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := GetUserStats(userID)
	if err != nil {
		logrus.Errorf("GetStats failed: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch stats: %v", err)})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func GetHistoricalINRHandler(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	history, err := GetHistoricalPortfolio(userID)
	if err != nil {
		logrus.Errorf("GetHistoricalINR failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func GetPortfolioHandler(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	portfolio, err := GetPortfolio(userID)
	if err != nil {
		logrus.Errorf("GetPortfolio failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolio"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}
