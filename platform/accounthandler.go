package platform

import (
	"strconv"

	coinbasepro "github.com/freddy212/go-coinbasepro"
)

func GetCurrencyPrice(accounts []coinbasepro.Account, currency string) float64 {
	for _, a := range accounts {
		if a.Currency == currency {
			amount, _ := strconv.ParseFloat(a.Balance, 64)
			return amount
		}
	}
	return 0.0
}

var buyTotal float64 = 0.0
var sellTotal float64 = 0.0

func BuyTotal(accounts []coinbasepro.Account, currency string) float64 {
	if buyTotal == 0.0 {
		buyTotal = GetCurrencyPrice(accounts, currency)
	}
	return buyTotal
}

func SellTotal(accounts []coinbasepro.Account, currency string) float64 {
	if sellTotal == 0.0 {
		sellTotal = GetCurrencyPrice(accounts, currency)
	}
	return sellTotal
}
