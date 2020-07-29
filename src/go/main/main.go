package main

import (
	backtest "TradingHackathon/src/go/backtestingUtils"
	live "TradingHackathon/src/go/liveUtils"
	. "TradingHackathon/src/go/rsi"
	"log"
	"os/exec"
	"time"

	"github.com/luno/luno-go/decimal"
)

// Global Variables
var isLive bool
var prevDay time.Time
var funds decimal.Decimal

func isNewDay() bool {
	y1, m1, d1 := prevDay.Date()
	y2, m2, d2 := time.Now().Date()
	return d1 != d2 || m1 != m2 || y1 != y2
}

func getPastAsks(b *RsiBot) []decimal.Decimal {
	//Populating past asks with 1 TradingPeriod worth of data
	pastAsks := make([]decimal.Decimal, b.TradingPeriod)
	var i int64 = 0
	for i < b.TradingPeriod {
		time.Sleep(20 * time.Second)
		pastAsks[i] = live.GetCurrAsk()
		i++
	}
	b.PrevAsk = pastAsks[b.TradingPeriod-1]
	return pastAsks
}

type tradeFunc func(b *RsiBot)

func main() {
	startBot("ETHXBT")
}

func startBot(pair string) {
	log.Println("Bot started:", pair)
	prevDay = time.Now().AddDate(0, 0, 0)

	// live.Email("START", decimal.Zero())

	isLive = true
	funds = decimal.NewFromInt64(100)
	var trade tradeFunc
	var pastAsks []decimal.Decimal

	live.PairName = pair
	live.Client = live.CreateClient()

	// initialising values within bot portfolio
	tradingPeriod := int64(14)
	StopLossMultDecimal := decimal.NewFromFloat64(0.9975, 8)
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
		UpEma:          decimal.Zero(),
		DownEma:        decimal.Zero(),
		PrevAsk:        decimal.Zero(),
	}

	log.Println("Getting past asks: STARTED")

	if isLive {
		trade = live.TradeLive
		pastAsks = getPastAsks(&bot)
	} else {
		backtest.InitialiseFunds(decimal.NewFromFloat64(0.014, 8), decimal.Zero())
		trade = backtest.TradeOffline

		var i int64
		for i = 0; i < tradingPeriod; i++ {
			pastAsks = append(pastAsks, backtest.GetOfflineAsk(i+1))
		}
	}

	log.Println("Getting past asks: COMPLETE")

	pastUps, pastDowns := []decimal.Decimal{}, []decimal.Decimal{}

	for i, v := range pastAsks {
		if i == 0 {
			continue
		}
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
		if isNewDay() {
			fileName := time.Now().Format("2006-01-02")
			live.ClosePrevFile(fileName)

			graphCmd := exec.Command("python3", "graphData.py", fileName)
			err1 := graphCmd.Run()

			if err1 != nil {
				log.Println("ERROR! Failed to graph data:", err1)
			}
			//Emailing
			//newFunds, _ := live.getAssets("XRP","XBT")
			newFunds := decimal.NewFromFloat64(0, 2)
			yield := newFunds.Sub(funds)
			live.Email("GRAPH", yield)
			funds = newFunds

			deletePicCmd := exec.Command("rm", "graph.png")
			err2 := deletePicCmd.Run()

			if err2 != nil {
				log.Println("ERROR! Failed to delete graph:", err2)
			}

			if err1 == nil && err2 == nil {
				log.Println("Graphed daily data successfully")
			}

			live.SetUpNewFile()
			bot.NumOfDecisions = 0
			prevDay = time.Now()
		}
		trade(&bot)
	}
}
