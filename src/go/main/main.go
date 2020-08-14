package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	live "../Utils"
	backtest "../backtestingUtils"
	"github.com/luno/luno-go/decimal"
)

// Global Variables
var prevDay time.Time
var funds decimal.Decimal


func isNewDay() bool {
	y1, m1, d1 := prevDay.Date()
	y2, m2, d2 := time.Now().Date()
	return d1 != d2 || m1 != m2 || y1 != y2
}

func GetCandlesticksandPastAsks(b *live.RsiBot) ([]live.Candlestick, []decimal.Decimal) {
	//Populating past asks with 1 MACDTradingPeriodLR worth of data
	var pastAsks []decimal.Decimal
	var stack []live.Candlestick
	var i int64 = 0
	for i < b.LongestTradingPeriod {
		stick := live.GetCandleStick(b.TimeInterval)
		pastAsks = append(pastAsks, stick.CloseAsk)
		status := "GETTING PAST ASKS [" + strconv.FormatInt(i+1, 10) + "/" + strconv.FormatInt((b.LongestTradingPeriod),10) + "]"
		live.PrintStatus(b, stick.CloseBid, stick.CloseAsk,status,"", nil)
		stack = append(stack, stick)
	//	fmt.Println("Candle " + strconv.Itoa(int(i)) + " got")
		i++
	}
	b.PrevAsk = pastAsks[b.LongestTradingPeriod-1]
	return stack, pastAsks
}

func findMax(array []int64) (max int64) {

	max = array[0]
	for _, value := range array {
		if value > max {
			max = value
		}
	}
	return max
}

type tradeFunc func(b *live.RsiBot)

func main() {
	startBot("ETHXBT")
}

func startBot(pair string) {

	live.LoadScreen()
	prevDay = time.Now().AddDate(0, 0, 0)
	live.StartDay = prevDay
	live.InitialiseKeys()


	funds = decimal.NewFromInt64(100)
	var trade tradeFunc
	var pastAsks []decimal.Decimal
	var stack []live.Candlestick

	live.PairName = pair
	live.User = strings.ToUpper(os.Args[1])
	live.Client = live.CreateClient()
	live.PrintStatus(nil, decimal.Zero(), decimal.Zero(), "UPDATE EMAIL SENT" ,"",nil)
	live.Email("START", decimal.Zero())

	timeInterval, _ := strconv.ParseInt(os.Args[3], 10, 64)

	botstring := ""
	botstring = os.Args[2]
	if botstring == "0000" {
		fmt.Println("No Strategies Chosen. Bot has been stopped")
		return
	}
	offset, _ := decimal.NewFromString("0.00000020")
	tradingperiodsused := []int64{}
	rsiTradingPeriod := int64(14)
	macdTradingPeriodLR := int64(10)
	macdTradingPeriodSR := int64(5)
	candleTradingPeriod := int64(3)
	offsetTradingPeriod := int64(14)
	if []rune(botstring)[0] != '0' {
		tradingperiodsused = append(tradingperiodsused, rsiTradingPeriod)
	}
	if []rune(botstring)[1] != '0' {
		tradingperiodsused = append(tradingperiodsused, macdTradingPeriodLR)
	}
	if []rune(botstring)[2] != '0' {
		tradingperiodsused = append(tradingperiodsused, candleTradingPeriod)
	}
	if []rune(botstring)[3] != '0' {
		tradingperiodsused = append(tradingperiodsused, offsetTradingPeriod)
	}

	longestTradingPeriod := findMax(tradingperiodsused)
	StopLossMultDecimal := decimal.NewFromFloat64(0.9975, 8)
	rsiLowerLim := decimal.NewFromInt64(20)

	pastAsks = []decimal.Decimal{}

	// initialising bot

	bot := live.RsiBot{
		RSITradingPeriod:     rsiTradingPeriod,
		MACDTradingPeriodLR:  macdTradingPeriodLR,
		MACDTradingPeriodSR:  macdTradingPeriodSR,
		LongestTradingPeriod: longestTradingPeriod,
		OffsetTraingPeriod:   offsetTradingPeriod,
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
		CandleTradingPeriod:  candleTradingPeriod,
		PastAsks:             pastAsks,
		TimeInterval:         timeInterval,
		Stack:                stack,
		Offset:               offset,
		BotString:            botstring,
	}

	//log.Println("User:", live.User)
	//log.Println("Getting past asks: STARTED")
	livecommand, _ := strconv.ParseInt(os.Args[4], 10, 64)
	if livecommand == 1 {
		trade = live.TradeLive
		bot.Stack, bot.PastAsks = GetCandlesticksandPastAsks(&bot)

	} else {
		backtest.InitialiseFunds(decimal.NewFromFloat64(0.014, 8), decimal.Zero())
		trade = backtest.TradeOffline

		var i int64
		for i = 0; i < longestTradingPeriod; i++ {
			bot.PastAsks = append(bot.PastAsks, backtest.GetOfflineAsk(i+1))
		}
	}
//	log.Println("Getting past asks: COMPLETE")
	pastUps, pastDowns := []decimal.Decimal{}, []decimal.Decimal{}
	for i, v := range bot.PastAsks[bot.LongestTradingPeriod-bot.RSITradingPeriod : bot.LongestTradingPeriod] {
		if i == 0 {
			continue
		}
		if v.Cmp(bot.PastAsks[i-1]) == -1 {
			pastDowns = append(pastDowns, bot.PastAsks[i-1].Sub(v))
		} else if v.Cmp(bot.PastAsks[i-1]) == 1 {
			pastUps = append(pastUps, v.Sub(bot.PastAsks[i-1]))
		}
	}

	bot.UpEma = live.InitialSma(pastUps, rsiTradingPeriod)
	bot.DownEma = live.InitialSma(pastDowns, rsiTradingPeriod)

	live.SetUpNewFile()
	for {
		if isNewDay() {
			fileName := time.Now().Format("2006-01-02")
			live.ClosePrevFile(fileName)

			graphCmd := exec.Command("python3", " graphData.py")
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
