package main

import (
	"fmt"
	"github.com/luno/luno-go/decimal"
	"time"
	backtest "TradingHackathon/src/go/backtestingUtils"
	live "TradingHackathon/src/go/liveUtils"
	. "TradingHackathon/src/go/rsi"
	""
)

// Global Variables
var isLive bool

func isMarketClosed() bool {
	start := "01:00"
	end := "01:05"
	check := time.Now()
  return !check.Before(start) && !check.After(end)
}



func getPastAsks(b *RsiBot) []decimal.Decimal {
	//Populating past asks with 1 TradingPeriod worth of data
	pastAsks := make([]decimal.Decimal, b.TradingPeriod)
	var i int64 = 0
	for i < b.TradingPeriod {
		time.Sleep(time.Minute)
		pastAsks[i] = live.GetCurrAsk()
		//delete from here to sleep
		buffer := ""
		if i < 9 {
			buffer = " "
		}

		fmt.Println("Filling past asks: ", buffer, i+1, "/", b.TradingPeriod, ":  BTC", pastAsks[i])
		i++
		//delete up to here
	}
	b.PrevAsk = pastAsks[b.TradingPeriod - 1]
	return pastAsks
}

type TradeFunc func(b *RsiBot)

func main() {

	live.Email("START", decimal.Zero(), decimal.Zero())

	isLive = true
	var trade TradeFunc
	var pastAsks []decimal.Decimal

	live.Pair = "XRPXBT"
	live.Client, live.ReqPointer = live.GetTickerRequest()
	live.Client.SetTimeout(time.Minute)

	// initialising values within bot portfolio
	tradingPeriod := int64(14)
	StopLossMultDecimal := decimal.NewFromFloat64(0.999, 8)
	rsiLowerLim := decimal.NewFromInt64(25)

	// initialising bot
	bot := RsiBot{
		TradingPeriod:  tradingPeriod,
		TradesMade:     0,
		NumOfDecisions: 0,
		StopLoss:       decimal.Zero(),
		StopLossMult:   StopLossMultDecimal,
		OverSold:       rsiLowerLim,
		ReadyToBuy:     true,
		BuyPrice:       decimal.Zero(),
		UpEma:					decimal.Zero(),
		DownEma:				decimal.Zero(),
		PrevAsk:				decimal.Zero(),
	}


	if isLive {
		trade = live.TradeLive
		pastAsks = getPastAsks(&bot)
	} else {
		backtest.InitialiseFunds(decimal.NewFromFloat64(0.014,8), decimal.Zero())
		trade = backtest.TradeOffline

		var i int64
		for i = 0; i < tradingPeriod; i++ {
			pastAsks = append(pastAsks, backtest.GetOfflineAsk(i+1))
		}
	}

	pastUps, pastDowns := []decimal.Decimal{}, []decimal.Decimal{}

	for i,v := range pastAsks {
		if i == 0 {continue}
		if v.Cmp(pastAsks[i-1]) == -1 {
			pastDowns = append(pastDowns, pastAsks[i-1].Sub(v))
		} else if v.Cmp(pastAsks[i-1]) == 1 {
			pastUps = append(pastUps, v.Sub(pastAsks[i-1]))
		}
	}

	bot.UpEma = Sma(pastUps, tradingPeriod)
	bot.DownEma = Sma(pastDowns, tradingPeriod)

	dataFile = excelize.NewFile()

	dataFile.SetCellValue("Sheet1", "A1", "Curr Price")
	dataFile.SetCellValue("Sheet1", "B1", "RSI")
	dataFile.SetCellValue("Sheet1", "C1", "Mode (ReadyToBuy?)")
	dataFile.SetCellValue("Sheet1", "D1", "BuyPrice")
	dataFile.SetCellValue("Sheet1", "D1", "SellPrice")
	for {
		if isMarketClosed(){

		}
		trade(&bot)
	}
}
