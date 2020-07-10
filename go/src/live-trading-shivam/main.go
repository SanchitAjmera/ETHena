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

func populatePastAsks(b *rsiBot) {
	//Populating past asks with 1 tradingPeriod worth of data
	var i int64 = 0
	for i < b.tradingPeriod {
		b.pastAsks[i] = getCurrAsk()
		//delete from here to sleep
		buffer := ""
		if i < 9 {
			buffer = " "
		}

		fmt.Println("Filling past asks: ", buffer, i+1, "/", b.tradingPeriod, ":  BTC", b.pastAsks[i])
		i++
		//delete up to here
		time.Sleep(time.Minute)
	}
}

// test function for the RSI bot
func test(b *rsiBot) {
	populatePastAsks(b)
	for {
		b.trade()
	}
}

func main() {

	pair = "XRPXBT"
	client, reqPointer = getTickerRequest()

	// initialising values within bot portfolio
	tradingPeriod := int64(15)
	stopLossMultDecimal := decimal.NewFromFloat64(0.995, 8)

	rsiLowerLim := decimal.NewFromInt64(30)

	// initialising bot
	bot := rsiBot{
		tradingPeriod:  tradingPeriod,
		tradesMade:     0,
		numOfDecisions: 0,
		stopLoss:       decimal.Zero(),
		stopLossMult:   stopLossMultDecimal,
		pastAsks:       make([]decimal.Decimal, tradingPeriod),
		overSold:       rsiLowerLim,
		readyToBuy:     true,
		buyPrice:       decimal.Zero(),
	}

	test(&bot)

}
