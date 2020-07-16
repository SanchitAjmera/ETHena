package main

import (
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"time"
)

// Global Variables
var client *luno.Client
var reqPointer *luno.GetTickerRequest
var pair string

func getPastAsks(b *rsiBot) []decimal.Decimal {
	//Populating past asks with 1 tradingPeriod worth of data
	pastAsks := make([]decimal.Decimal, b.tradingPeriod)
	var i int64 = 0
	for i < b.tradingPeriod {
		time.Sleep(time.Minute)
		pastAsks[i] = getCurrAsk()
		//delete from here to sleep
		buffer := ""
		if i < 9 {
			buffer = " "
		}

		fmt.Println("Filling past asks: ", buffer, i+1, "/", b.tradingPeriod, ":  BTC", pastAsks[i])
		i++
		//delete up to here
	}
	b.prevAsk = pastAsks[b.tradingPeriod - 1]
	return pastAsks
}

// test function for the RSI bot
func test(b *rsiBot) {
	pastAsks := getPastAsks(b)
	b.upEma = sma(pastAsks, b.tradingPeriod)
	b.downEma = sma(pastAsks, b.tradingPeriod)
	for {
		b.trade()
	}
}

func main() {

	pair = "XRPXBT"
	client, reqPointer = getTickerRequest()
	client.SetTimeout(time.Minute)

	// initialising values within bot portfolio
	tradingPeriod := int64(14)
	stopLossMultDecimal := decimal.NewFromFloat64(0.999, 8)
	rsiLowerLim := decimal.NewFromInt64(25)

	// initialising bot
	bot := rsiBot{
		tradingPeriod:  tradingPeriod,
		tradesMade:     0,
		numOfDecisions: 0,
		stopLoss:       decimal.Zero(),
		stopLossMult:   stopLossMultDecimal,
		overSold:       rsiLowerLim,
		readyToBuy:     true,
		buyPrice:       decimal.Zero(),
		upEma:					decimal.Zero(),
		downEma:				decimal.Zero(),
		prevAsk:				decimal.Zero(),
	}

	test(&bot)

}
