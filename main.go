package main

import (
	"time"

	"github.com/frede/gocustomtrading/platform"
)

func main() {
	var listenedCurrencies = []string{
		"ADA-EUR",
	}
	var decimalsToSell = []string{"1"}

	wait := make(chan bool)

	start := time.Now()
	index := 0

	for {
		newTime := time.Now()
		duration, _ := time.ParseDuration("2s")
		if newTime.Sub(start) > duration && index < len(listenedCurrencies) {
			go platform.StartSocket(listenedCurrencies[index], decimalsToSell[index])
			start = time.Now()
			index += 1
		}
		if index > len(decimalsToSell)-1 {
			{
				<-wait
				println("not looping")
			}
		}
	}
}
