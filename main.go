package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// 1. Config
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Error loading .env file, using system environment variables")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	serverPort := os.Getenv("PORT")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || serverPort == "" {
		logrus.Fatal("Missing required environment variables")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// 2. Database
	ConnectDB(dsn)

	// 3. Start Background Worker
	StartPriceUpdater()

	// 4. Router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api")
	{
		api.POST("/reward", CreateRewardHandler)
		api.GET("/today-stocks/:userId", GetTodayStocksHandler)
		api.GET("/historical-inr/:userId", GetHistoricalINRHandler)
		api.GET("/stats/:userId", GetStatsHandler)
		api.GET("/portfolio/:userId", GetPortfolioHandler)
	}

	// 5. Run Server
	logrus.Infof("Server starting on port %s", serverPort)
	if err := r.Run(":" + serverPort); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
