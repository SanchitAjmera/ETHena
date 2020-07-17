package main

import (
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"time"
	backtest "TradingHackathon/src/go/backtestingUtils"
	live "TradingHackathon/src/go/liveUtils"
)

// Global Variables
var isLive bool

// struct for the rsiBot
type rsiBot struct {
	tradingPeriod  int64             // No of past asks used to calculate RSI
	tradesMade     int64             // total number of trades executed
	numOfDecisions int64             // number of times the bot calculates
	stopLoss       decimal.Decimal   // variable stop loss
	stopLossMult   decimal.Decimal   // multiplier for stop loss
	overSold       decimal.Decimal   // bound to tell the bot when to buy
	readyToBuy     bool              // false means ready to sell
	buyPrice       decimal.Decimal   // stores most recent price we bought at
	upEma					 decimal.Decimal   // exponentially smoothed Wilder's MMA for upward change
	downEma 			 decimal.Decimal   // exponentially smoothed Wilder's MMA for downward change
	prevAsk				 decimal.Decimal	 // the previous recorded ask price
}

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

type TradeFunc func(b *rsiBot)

func main() {

	isLive = false
	var TradeFunc trade

	if isLive {
		trade = live.tradeLive
	} else {
		initialiseFunds()
		trade = backtest.tradeOffline
	}

	if isLive()

	live.Pair = "XRPXBT"
	live.Client, live.ReqPointer = getTickerRequest()
	live.Client.SetTimeout(time.Minute)

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
