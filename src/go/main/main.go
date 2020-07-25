package main

import (
	"fmt"
	"github.com/luno/luno-go/decimal"
	"time"
	backtest "TradingHackathon/src/go/backtestingUtils"
	live "TradingHackathon/src/go/liveUtils"
	. "TradingHackathon/src/go/rsi"
	"os/exec"
)

// Global Variables
var isLive bool
var prevDay time.Time

func isNewDay() bool {
    y1, m1, d1 := prevDay.Date()
    y2, m2, d2 := time.Now().Date()
    return  d1 != d2 ||  m1 != m2 || y1 != y2
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

	prevDay = time.Now().AddDate(0, 0, -1)

	live.Email("START", decimal.Zero())

	isLive = true
	var trade TradeFunc
	var pastAsks []decimal.Decimal

	live.Pair = "XRPXBT"
	live.Client, live.ReqPointer = live.GetTickerRequest()
	live.Client.SetTimeout(time.Minute)

	// initialising values within bot portfolio
	tradingPeriod := int64(14)
	StopLossMultDecimal := decimal.NewFromFloat64(0.999, 8)
	rsiLowerLim := decimal.NewFromInt64(20)

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

	live.SetUpNewFile()
	for {
		if isNewDay(){
			fileName := time.Now().Format("2006-01-02")
			live.ClosePrevFile(fileName)

			graphCmd := exec.Command("python3","graphData.py", fileName)
			err1 := graphCmd.Run()

			if err1 != nil {
				fmt.Println("ERROR! Failed to graph data:", err1)
			}

			//Emailing

			deletePicCmd := exec.Command("rm", "graph.png")
			err2 := deletePicCmd.Run()

			if err2 != nil {
				fmt.Println("ERROR! Failed to delete graph:", err2)
			}

			if err1 == nil && err2 == nil {
				fmt.Println("Graphed daily data successfully")
			}

			live.SetUpNewFile()
			bot.NumOfDecisions = 0
			prevDay = time.Now()
		}
		trade(&bot)
	}
}
