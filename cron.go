package main

//test 1
import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

func StartPriceUpdater() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		UpdatePrices() // Run immediately on start
		for {
			select {
			case <-ticker.C:
				UpdatePrices()
			}
		}
	}()
}

func UpdatePrices() {
	logrus.Info("Running hourly price update...")
	stocks := []string{"RELIANCE", "TCS", "INFY", "HDFCBANK"}

	for _, symbol := range stocks {
		price, _ := GetCurrentPrice(symbol)
		err := InsertStockPrice(symbol, price)
		if err != nil {
			logrus.Errorf("Failed to update price for %s: %v", symbol, err)
		}
	}
	logrus.Info("Price update completed.")
}

func GetCurrentPrice(symbol string) (float64, error) {
	rand.Seed(time.Now().UnixNano())
	return 1000 + rand.Float64()*(3000-1000), nil
}
