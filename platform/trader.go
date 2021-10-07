package platform

import (
	"fmt"
	"sort"

	coinbasepro "github.com/freddy212/go-coinbasepro"
)

func Sell(coin string, productId string, client *coinbasepro.Client, decimalToSell string, price float64, tradeHistory *TradeHistory) {
	tradeHistory.SellStreak = min(tradeHistory.SellStreak+1, 8)
	priceAt := findIndex(tradeHistory.TradePrices, price)
	sellMultipler := getSellMultiplier(len(tradeHistory.TradePrices), priceAt)

	sellPercent := 10 * sellMultipler * (0.5 + float64(tradeHistory.SellStreak)/2)

	sellAmount := tradeHistory.sellTotal / 100 * sellPercent
	println("Attempting to sell " + coin)
	println("decimal to sell is :" + decimalToSell)
	println("amount to sell is :", fmt.Sprintf("%."+decimalToSell+"f", sellAmount))
	println("SellStreak is ", tradeHistory.SellStreak)
	println("SellMultiplier is ", sellMultipler)
	order := coinbasepro.Order{
		Type:      "market",
		Size:      fmt.Sprintf("%."+decimalToSell+"f", sellAmount),
		Side:      "sell",
		ProductID: productId,
	}
	savedOrder, err := client.CreateOrder(&order)
	if err != nil {
		println(err.Error())
	}
	println(savedOrder.ID)
	tradeHistory.TradePrices = append(tradeHistory.TradePrices, price)
	sort.Float64s(tradeHistory.TradePrices)
	tradeHistory.BuyStreak = 0
}
func Buy(productId string, client *coinbasepro.Client, price float64, tradeHistory *TradeHistory) {
	tradeHistory.BuyStreak = min(tradeHistory.BuyStreak+1, 8)
	priceAt := findIndex(tradeHistory.TradePrices, price)
	buyMultipler := getBuyMultiplier(len(tradeHistory.TradePrices), priceAt)

	buyPercent := 10 * buyMultipler * (0.5 + float64(tradeHistory.BuyStreak)/2)
	buyAmount := tradeHistory.buyTotal / 100 * buyPercent

	println("Attempting to buy " + productId)
	println("Amount to buy is ", buyAmount)
	println("BuyStreak is ", tradeHistory.BuyStreak)
	println("BuyMultiplier is ", buyMultipler)

	order := coinbasepro.Order{
		Type:      "market",
		Funds:     fmt.Sprintf("%.2f", buyAmount),
		Side:      "buy",
		ProductID: productId,
	}

	savedOrder, err := client.CreateOrder(&order)
	if err != nil {
		println(err.Error())
	}
	println(savedOrder.ID)
	tradeHistory.TradePrices = append(tradeHistory.TradePrices, price)
	sort.Float64s(tradeHistory.TradePrices)
	tradeHistory.SellStreak = 0
}

func findIndex(prices []float64, newPrice float64) int {
	startValue := 0
	for i, price := range prices {
		if newPrice < price {
			break
		}
		startValue = i
	}
	println("bigger price found at", startValue)
	fmt.Printf("%v", prices)
	return startValue
}

func getBuyMultiplier(priceHistorySize int, priceIndex int) float64 {
	if priceIndex == 0 {
		return 2
	} else if priceIndex <= priceHistorySize/2 {
		return 1.5
	} else {
		return 1
	}
}
func getSellMultiplier(priceHistorySize int, priceIndex int) float64 {
	if priceIndex == priceHistorySize-1 {
		return 2
	} else if priceIndex >= priceHistorySize/2 {
		return 1.5
	} else {
		return 1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
