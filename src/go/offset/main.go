package main

import (
	"fmt"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
  backtest "TradingHackathon/src/go/backtestingUtils"
  "time"

)

// Global Variables
var isLive bool
var Client *luno.Client
var ReqPointer *luno.GetTickerRequest
var Pair string
var currRow int64

func getPastAsks(b *offsetBot) []decimal.Decimal {
	//Populating past asks with 1 tradingPeriod worth of data
	pastAsks := make([]decimal.Decimal, b.tradingPeriod)
	var i int64 = 0
	for i < b.tradingPeriod {
		time.Sleep(time.Minute)
		pastAsks[i] = GetCurrAsk()
		//delete from here to sleep
		buffer := ""
		if i < 9 {
			buffer = " "
		}

		fmt.Println("Filling past asks: ", buffer, i+1, "/", b.tradingPeriod, ":  BTC", pastAsks[i])
		i++
		//delete up to here
	}
	b.PrevAsk = pastAsks[b.tradingPeriod-1]
	return pastAsks
}

func main() {

  isLive = true

	var pastAsks []decimal.Decimal

	Pair = "XRPXBT"
	Client, ReqPointer = GetTickerRequest()
	Client.SetTimeout(time.Minute)

	// initialising values within bot portfolio
	tradingPeriod := int64(14)
	currRow = tradingPeriod + 2
	StopLossMultDecimal := decimal.NewFromFloat64(1, 8)
	offset, _ := decimal.NewFromString("0.00000020")

	// initialising bot
	offsetBot := offsetBot{
		tradingPeriod: tradingPeriod,  // How often the bot calculates a long term result
		ema:           decimal.Zero(), // exponentially smoothed Wilder's MMA for upward change
		offset:        offset,
		readyToBuy:    true,
		StopLoss:      decimal.Zero(),
		StopLossMult:  StopLossMultDecimal,
		BuyPrice:      decimal.Zero(),
		PrevAsk:       decimal.Zero(),
	}

	if isLive {
		email("START", decimal.Zero(), decimal.Zero())
		pastAsks = getPastAsks(&offsetBot)
		offsetBot.ema = sma(pastAsks)
		for {
			offsetBot.tradeOnline()
		}

	} else {

	  backtest.InitialiseFunds(decimal.NewFromFloat64(0.014,8), decimal.Zero())
		for i:= 0; i < int(tradingPeriod); i++ {
			pastAsks = append(pastAsks, backtest.GetOfflineAsk(int64(i+1)))

		}
		fmt.Println("\n\n\n OFFSET:", offsetBot.offset, "\n")
		offsetBot.ema = sma(pastAsks)
		for i := 0; i < 1300; i++ {
			offsetBot.tradeOffline()
		}
	}

}
