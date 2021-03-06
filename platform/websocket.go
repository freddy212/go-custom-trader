package platform

import (
	"strconv"
	"strings"

	coinbasepro "github.com/freddy212/go-coinbasepro"
	ws "github.com/gorilla/websocket"
)

type Counter struct {
	TickCount       int
	PriceTotal      float64
	Angle           float64
	PreviousPrice   float64
	Average         float64
	AngleTick       []float64
	LongTermPrice   float64
	LongTermAverage float64
	DecimalToSell   string
}

type TradeHistory struct {
	TradePrices []float64
	SellStreak  int
	BuyStreak   int
	buyTotal    float64
	sellTotal   float64
}

func StartSocket(productId string, decimalCount string) {

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
	if err != nil {
		println(err.Error())
	}

	coinName := strings.Split(productId, "-")[0]

	subscribe := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{
			coinbasepro.MessageChannel{
				Name:       "heartbeat",
				ProductIds: []string{productId},
			},
			coinbasepro.MessageChannel{
				Name:       "ticker",
				ProductIds: []string{productId},
			},
		},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}
	client := GetClientInstance()
	var counter Counter
	var tradeHistory TradeHistory
	accounts, _ := client.GetAccounts()
	tradeHistory.buyTotal = BuyTotal(accounts, "EUR")

	println("buy total is :", tradeHistory.buyTotal)
	counter.DecimalToSell = decimalCount
	println("started listening for: " + productId)
	println("decimal to sell is : " + counter.DecimalToSell)
	for {
		message := coinbasepro.Message{}
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			println("trying to continue")
			wsConn, _, _ = wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
			if err := wsConn.WriteJSON(subscribe); err != nil {
				println(err.Error())
			}
		}
		if message.Type == "ticker" {
			price, _ := strconv.ParseFloat(message.Price, 64)
			analyzePrice(&counter, price, client, coinName, message.ProductID, &tradeHistory)
		}
	}
}
func analyzePrice(counter *Counter, price float64, client *coinbasepro.Client, coinName string, productId string, tradeHistory *TradeHistory) {
	counter.PriceTotal += price
	counter.TickCount++

	if counter.TickCount > 1 {
		counter.Angle += (price - counter.PreviousPrice)
		counter.AngleTick = append(counter.AngleTick, (price - counter.PreviousPrice))
	}
	if counter.TickCount%50 == 0 && counter.Average == 0.0 {
		counter.Average = counter.PriceTotal / 50
		counter.PriceTotal = 0.0
		println("average is ", counter.Average)
		println("for coin " + coinName)
	}
	if counter.TickCount%30 == 0 {
		println("tick for coin ", coinName)
		println("price is", price)
		println("angle is", counter.Angle)
	}
	if counter.TickCount > 151 {
		counter.Angle -= counter.AngleTick[0]
		//counter.AngleTick = append(counter.AngleTick[:0], counter.AngleTick[1:]...)
		counter.AngleTick = counter.AngleTick[1:]
	}
	counter.PreviousPrice = price
	if counter.Average != 0.0 {
		if price > UpperBound(counter.Average) && counter.Angle < 0 {
			println("average before sell : %e", counter.Average)
			println("price is: ", price)
			tradeHistory.sellTotal = 50 / price
			println("sell total for coin ", coinName)
			println("is", tradeHistory.sellTotal)

			Sell(coinName, productId, client, counter.DecimalToSell, price, tradeHistory)
			counter.Average = price
			counter.Angle = 0
			counter.PriceTotal = 0.0
			counter.TickCount = 0
			counter.AngleTick = nil

		}
		if price < LowerBound(counter.Average) && counter.Angle > 0 {
			println("average before buy : %e", counter.Average)
			println("price is: ", price)
			Buy(productId, client, price, tradeHistory)
			counter.Average = price
			counter.Angle = 0.0
			counter.PriceTotal = 0.0
			counter.TickCount = 0
			counter.AngleTick = nil
		}
	}
}

func UpperBound(num float64) float64 {
	return num / 100 * 102.5
}
func LowerBound(num float64) float64 {
	return num / 100 * 97.5
}
