package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	backtest "../backtestingUtils"
	live "../liveUtils"
	. "../rsi"
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
	//Populating past asks with 1 MACDTradingPeriodLR worth of data
	pastAsks := make([]decimal.Decimal, b.LongestTradingPeriod)
	var i int64 = 0
	for i < b.LongestTradingPeriod {
		time.Sleep(live.TimeDuration * time.Second)
		pastAsks[i] = live.GetCurrAsk()
		i++
	}
	b.PrevAsk = pastAsks[b.LongestTradingPeriod-1]
	return pastAsks
}

func Maxof2int(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

type tradeFunc func(b *RsiBot)

func main() {
	startBot("ETHXBT")
}

func startBot(pair string) {
	log.Println("Bot started:", pair)
	prevDay = time.Now().AddDate(0, 0, 0)
	live.InitialiseKeys()

	// live.Email("START", decimal.Zero())

	isLive = false

	funds = decimal.NewFromInt64(100)
	var trade tradeFunc
	var pastAsks []decimal.Decimal

	live.PairName = pair
	live.User = strings.ToUpper(os.Args[1])
	live.Client = live.CreateClient()
	live.VOLUME_TIME_PERIOD = 5
	live.PROFIT_TIME_PERIOD = 30
	if os.Args[2] == "volume" {
		live.TimeDuration = live.VOLUME_TIME_PERIOD
	} else if os.Args[2] == "profit" {
		live.TimeDuration = live.PROFIT_TIME_PERIOD
	}

	// initialising values within bot portfolio
	rsiTradingPeriod := int64(14)
	macdTradingPeriodLR := int64(60)
	macdTradingPeriodSR := int64(30)
	longestTradingPeriod := int64(0)
	longestTradingPeriod = Maxof2int(Maxof2int(rsiTradingPeriod, macdTradingPeriodSR), macdTradingPeriodLR)
	StopLossMultDecimal := decimal.NewFromFloat64(0.9975, 8)
	rsiLowerLim := decimal.NewFromInt64(20)

	pastAsks = []decimal.Decimal{}

	log.Println("Getting past asks: COMPLETE")
	// initialising bot

	bot := RsiBot{
		RSITradingPeriod:     rsiTradingPeriod,
		MACDTradingPeriodLR:  macdTradingPeriodLR,
		MACDTradingPeriodSR:  macdTradingPeriodSR,
		LongestTradingPeriod: longestTradingPeriod,
		TradesMade:           0,
		NumOfDecisions:       0,
		StopLoss:             decimal.Zero(),
		StopLossMult:         StopLossMultDecimal,
		OverSold:             rsiLowerLim,
		ReadyToBuy:           true,
		BuyPrice:             decimal.Zero(),
		UpEma:                decimal.Zero(),
		DownEma:              decimal.Zero(),
		PrevAsk:              decimal.Zero(),
		MACDlongperiodavg:    decimal.Zero(),
		MACDshortperiodavg:   decimal.Zero(),
		PastAsks:             pastAsks,
	}

	log.Println("User:", live.User)
	log.Println("Getting past asks: STARTED")
	if isLive {
		trade = live.TradeLive
		pastAsks = getPastAsks(&bot)
	} else {
		backtest.InitialiseFunds(decimal.NewFromFloat64(0.014, 8), decimal.Zero())
		trade = backtest.TradeOffline

		var i int64
		for i = 0; i < longestTradingPeriod; i++ {
			pastAsks = append(pastAsks, backtest.GetOfflineAsk(i+1))
		}
	}
	bot.PastAsks = pastAsks
	pastUps, pastDowns := []decimal.Decimal{}, []decimal.Decimal{}
	for i, v := range pastAsks[longestTradingPeriod-rsiTradingPeriod : longestTradingPeriod] {
		if i == 0 {
			continue
		}
		if v.Cmp(pastAsks[i-1]) == -1 {
			pastDowns = append(pastDowns, pastAsks[i-1].Sub(v))
		} else if v.Cmp(pastAsks[i-1]) == 1 {
			pastUps = append(pastUps, v.Sub(pastAsks[i-1]))
		}
	}

	bot.UpEma = InitialSma(pastUps, rsiTradingPeriod)
	bot.DownEma = InitialSma(pastDowns, rsiTradingPeriod)

	live.SetUpNewFile()
	for {
		if isNewDay() {
			fileName := time.Now().Format("2006-01-02")
			live.ClosePrevFile(fileName)

			graphCmd := exec.Command("python3", "graphData.py")
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
